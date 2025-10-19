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

	"mahaam-api/handler"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/configs"
	"mahaam-api/internal/pkg/dbs"
	logs "mahaam-api/internal/pkg/log"
	"mahaam-api/internal/pkg/middleware"
	"mahaam-api/internal/pkg/monitor"
	"mahaam-api/internal/user"

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
	configs.InitConfigs()
	user.InitEmail()
	db, err := dbs.Init()
	if err != nil {
		logs.Error(uuid.Nil, "Failed to connect to database: %v", err)
	}
	defer db.Close()

	logs.Init(monitor.CreateLog)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authed := router.Group("/mahaam-api")
	authed.Use(middleware.TrafficMiddleware())
	authed.Use(middleware.RecoveryMiddleware())
	authed.Use(middleware.AuthMiddleware())

	// Register routes
	handler.RegisterAudit(authed)
	handler.RegisterHealth(authed)
	handler.RegisterUsers(authed)
	handler.RegisterPlans(authed)
	handler.RegisterTasks(authed)

	// Initialize health monitoring
	nodeName, _ := os.Hostname()
	health := &model.Health{
		ID:         uuid.New(),
		ApiName:    configs.ApiName,
		ApiVersion: configs.ApiVersion,
		NodeIP:     getNodeIP(),
		NodeName:   nodeName,
		EnvName:    configs.EnvName,
	}
	fmt.Println("health", health)
	monitor.ServerStarted(health)
	configs.InitCache(health)

	startMsg := fmt.Sprintf("✓ %s-v%s/%s-%s started with healthID=%s", configs.ApiName, configs.ApiVersion, configs.NodeIP, configs.NodeName, configs.HealthID.String())
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
	monitor.StartSendingPulses(pulseCtx, configs.HealthID)

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

	monitor.ServerStopped(configs.HealthID)
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(uuid.Nil, "Server forced to shutdown: %v", err)
	}
}
