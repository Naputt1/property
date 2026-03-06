package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/internal/services"
	"backend/internal/services/queue"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
)

// @title UK Housing Market & Rental Intelligence API
// @version 1.0
// @description This is a comprehensive API for analyzing and accessing UK housing market and rental data.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @query.collection.format multi

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization

func main() {
	flag.Parse()

	// Initialize Configuration & Logging
	cfg, err := config.CreateConfig()
	if err != nil {
		slog.Error("failed to create config", "error", err)
		os.Exit(1)
	}
	defer cfg.Close()

	// Initialize Database Connection and auto-migrations
	db.InitDB(cfg)

	// Initialize asynq client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", cfg.Opt.Redis.Host, cfg.Opt.Redis.Port),
		Password: cfg.Opt.Redis.Password,
		DB:       cfg.Opt.Redis.DB,
	})
	defer asynqClient.Close()

	// Initialize Repositories
	repos := services.Repositories{
		User:      repository.NewUserRepository(cfg.DB),
		Property:  repository.NewPropertyRepository(cfg.DB),
		Job:       repository.NewJobRepository(cfg.DB),
		Analytics: repository.NewAnalyticsRepository(cfg.DB),
	}

	// Initialize Services
	svcs := services.NewServices(cfg, repos, asynqClient)

	// Initialize and Start Asynq Server for background migration tasks
	asynqSrv, err := queue.NewAsynqServer(cfg, svcs)
	if err != nil {
		slog.Error("failed to initialize asynq server", "error", err)
		os.Exit(1)
	}
	
	go func() {
		if err := asynqSrv.Start(); err != nil {
			slog.Error("Failed to start asynq server", "error", err)
		}
	}()
	defer asynqSrv.Stop()

	// Initialize HTTP Router
	r := routes.SetupRouter(cfg, svcs)

	port := cfg.Opt.Port
	if port == "" {
		port = "8080"
	}

	rest := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start HTTP Server
	go func() {
		slog.Info("Starting Property API Server", "port", port)
		if err := rest.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := rest.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
