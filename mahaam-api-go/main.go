package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mahaam-api/feat/handler"
	"mahaam-api/feat/models"
	"mahaam-api/feat/repo"
	"mahaam-api/feat/service"
	"mahaam-api/infra/cache"
	"mahaam-api/infra/configs"
	"mahaam-api/infra/dbs"
	"mahaam-api/infra/emails"
	logs "mahaam-api/infra/log"
	"mahaam-api/infra/middleware"
	"mahaam-api/infra/security"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func getNodeIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:10002")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	configs.Init()
	emails.Init()
	db, err := dbs.Init()
	if err != nil {
		logs.Error(uuid.Nil, "Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Init repos
	planRepo := repo.NewPlanRepo()
	planMembersRepo := repo.NewPlanMembersRepo()
	userRepo := repo.NewUserRepo()
	suggestedEmailsRepo := repo.NewSuggestedEmailRepo()
	taskRepo := repo.NewTaskRepo()
	deviceRepo := repo.NewDeviceRepo()
	logRepo := repo.NewLogRepo()
	trafficRepo := repo.NewTrafficRepo()
	healthRepo := repo.NewHealthRepo()

	logs.Init(logRepo.Create)

	// infraa
	authService := security.NewAuthService(deviceRepo, userRepo)

	// Init services
	healthService := service.NewHealthService(healthRepo)
	planService := service.NewPlanService(db, planRepo, planMembersRepo, userRepo, suggestedEmailsRepo)
	taskService := service.NewTaskService(taskRepo, planRepo)
	userService := service.NewUserService(db,
		userRepo,
		deviceRepo,
		planRepo,
		suggestedEmailsRepo,
		authService,
	)

	// Init handlers
	userHandler := handler.NewUserHandler(userService)
	planHandler := handler.NewPlanHandler(planService)
	auditHandler := handler.NewAuditHandler(logRepo)
	healthHandler := handler.NewHealthHandler()
	taskHandler := handler.NewTaskHandler(taskService)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authed := router.Group("/mahaam-api")
	authed.Use(middleware.TrafficMiddleware(trafficRepo))
	authed.Use(middleware.RecoveryMiddleware())
	authed.Use(middleware.AuthMiddleware(&authService, logRepo))

	// Register routes
	handler.RegisterUserHandler(authed, userHandler)
	handler.RegisterPlanHandler(authed, planHandler)
	handler.RegisterTaskHandler(authed, taskHandler)
	handler.RegisterAuditHandler(authed, auditHandler)
	handler.RegisterHealthHandler(authed, healthHandler)

	// Initialize health monitoring
	nodeName, _ := os.Hostname()
	health := &models.Health{
		ID:         uuid.New(),
		ApiName:    configs.ApiName,
		ApiVersion: configs.ApiVersion,
		NodeIP:     getNodeIP(),
		NodeName:   nodeName,
		EnvName:    configs.EnvName,
	}
	fmt.Println("health", health)
	healthService.ServerStarted(health)
	cache.Init(health)

	startMsg := fmt.Sprintf("✓ %s-v%s/%s-%s started with healthID=%s", cache.ApiName, cache.ApiVersion, cache.NodeIP, cache.NodeName, cache.HealthID.String())
	logs.Info(uuid.Nil, startMsg)
	time.Sleep(2 * time.Second)

	// Start the server
	port := fmt.Sprintf(":%v", configs.HTTPPort)
	fmt.Printf("Server running on port %s\n", port)

	srv := &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start pulse sending
	pulseCtx, pulseCancel := context.WithCancel(context.Background())
	healthService.StartSendingPulses(pulseCtx)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Error(uuid.Nil, "Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop pulse sending
	pulseCancel()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthService.ServerStopped()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(uuid.Nil, "Server forced to shutdown: %v", err)
	}
}
