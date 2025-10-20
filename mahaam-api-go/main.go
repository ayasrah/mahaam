package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mahaam-api/app/handler"
	"mahaam-api/app/models"
	"mahaam-api/app/repo"
	"mahaam-api/app/service"
	"mahaam-api/utils/conf"
	emails "mahaam-api/utils/email"
	logs "mahaam-api/utils/log"
	"mahaam-api/utils/middleware"
	token "mahaam-api/utils/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type repos struct {
	plan            repo.PlanRepo
	planMembers     repo.PlanMembersRepo
	user            repo.UserRepo
	suggestedEmails repo.SuggestedEmailRepo
	task            repo.TaskRepo
	device          repo.DeviceRepo
	log             repo.LogRepo
	traffic         repo.TrafficRepo
	health          repo.HealthRepo
}

type services struct {
	health service.HealthService
	plan   service.PlanService
	task   service.TaskService
	user   service.UserService
}

type handlers struct {
	user   handler.UserHandler
	plan   handler.PlanHandler
	audit  handler.AuditHandler
	health handler.HealthHandler
	task   handler.TaskHandler
}

func loadConfig() *conf.Conf {
	return conf.NewConf("config.json")
}

func openDB(dbURL string) *repo.AppDB {
	db, err := repo.NewAppDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	return db
}

func initRepos(db *repo.AppDB) repos {
	return repos{
		plan:            repo.NewPlanRepo(db),
		planMembers:     repo.NewPlanMembersRepo(db),
		user:            repo.NewUserRepo(db),
		suggestedEmails: repo.NewSuggestedEmailRepo(db),
		task:            repo.NewTaskRepo(db),
		device:          repo.NewDeviceRepo(db),
		log:             repo.NewLogRepo(db),
		traffic:         repo.NewTrafficRepo(db),
		health:          repo.NewHealthRepo(db),
	}
}

func initUtilities(cfg *conf.Conf, logRepo repo.LogRepo, deviceRepo repo.DeviceRepo, userRepo repo.UserRepo) (logs.Logger, token.TokenService, emails.EmailService) {
	logger := logs.NewLogger(cfg, logRepo.Create)
	tokenService := token.NewTokenService(deviceRepo, userRepo, cfg)
	emailService := emails.NewEmailService(cfg, logger)
	return logger, tokenService, emailService
}

func initServices(cfg *conf.Conf, logger logs.Logger, db *repo.AppDB, r repos, tokenService token.TokenService, emailService emails.EmailService) services {
	return services{
		health: service.NewHealthService(r.health, cfg, logger),
		plan:   service.NewPlanService(db, r.plan, r.planMembers, r.user, r.suggestedEmails),
		task:   service.NewTaskService(db, r.task, r.plan),
		user:   service.NewUserService(db, r.user, r.device, r.plan, r.suggestedEmails, tokenService, emailService, cfg, logger),
	}
}

func initHandlers(svcs services, logger logs.Logger, cfg *conf.Conf) handlers {
	return handlers{
		user:   handler.NewUserHandler(svcs.user, logger),
		plan:   handler.NewPlanHandler(svcs.plan, logger),
		audit:  handler.NewAuditHandler(logger),
		health: handler.NewHealthHandler(cfg),
		task:   handler.NewTaskHandler(svcs.task),
	}
}

func buildRouter(cfg *conf.Conf, r repos, logger logs.Logger, tokenService token.TokenService, h handlers) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authed := router.Group("/mahaam-api")
	authed.Use(middleware.TrafficMiddleware(r.traffic, cfg, logger))
	authed.Use(middleware.RecoveryMiddleware(logger))
	authed.Use(middleware.AuthMiddleware(&tokenService, logger))

	// Register routes
	handler.RegisterUserHandler(authed, h.user)
	handler.RegisterPlanHandler(authed, h.plan)
	handler.RegisterTaskHandler(authed, h.task)
	handler.RegisterAuditHandler(authed, h.audit)
	handler.RegisterHealthHandler(authed, h.health)

	return router
}

func getNodeIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:10002")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func initHealthState(cfg *conf.Conf, healthSvc service.HealthService, logger logs.Logger) {
	nodeName, _ := os.Hostname()
	health := &models.Health{
		ID:         uuid.New(),
		ApiName:    cfg.ApiName,
		ApiVersion: cfg.ApiVersion,
		NodeIP:     getNodeIP(),
		NodeName:   nodeName,
		EnvName:    cfg.EnvName,
	}

	healthSvc.ServerStarted(health)
	conf.NewEnvironment(health)

	startMsg := fmt.Sprintf("✓ %s-v%s/%s-%s started with healthID=%s", cfg.ApiName, cfg.ApiVersion, conf.Env().NodeIP, conf.Env().NodeName, conf.Env().HealthID.String())
	logger.Info(uuid.Nil, startMsg)
	time.Sleep(2 * time.Second)
}

func startHTTPServer(router *gin.Engine, port int) *http.Server {
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("Server running on port %s\n", addr)

	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	return srv
}

func startPulse(healthSvc service.HealthService) (context.Context, context.CancelFunc) {
	pulseCtx, pulseCancel := context.WithCancel(context.Background())
	healthSvc.StartSendingPulses(pulseCtx)
	return pulseCtx, pulseCancel
}

func gracefulShutdown(srv *http.Server, healthSvc service.HealthService, logger logs.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthSvc.ServerStopped()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(uuid.Nil, "Server forced to shutdown: %v", err)
	}
}

func main() {
	cfg := loadConfig()
	db := openDB(cfg.DBUrl)
	defer db.Close()

	r := initRepos(db)
	logger, tokenService, emailService := initUtilities(cfg, r.log, r.device, r.user)
	svcs := initServices(cfg, logger, db, r, tokenService, emailService)
	h := initHandlers(svcs, logger, cfg)

	router := buildRouter(cfg, r, logger, tokenService, h)

	initHealthState(cfg, svcs.health, logger)

	_, pulseCancel := startPulse(svcs.health)
	defer pulseCancel()

	srv := startHTTPServer(router, cfg.HTTPPort)

	gracefulShutdown(srv, svcs.health, logger)
}
