package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/financial-analytics/api-gateway/internal/config"
    "github.com/financial-analytics/api-gateway/internal/gateway"
    "github.com/financial-analytics/api-gateway/internal/middleware"
    "github.com/financial-analytics/api-gateway/internal/services"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        logger.Fatal("Failed to load configuration", zap.Error(err))
    }

    // Initialize services
    authService := services.NewAuthService(cfg)
    rateLimiter := middleware.NewRateLimiter(cfg.Redis)
    
    // Create gateway
    gw := gateway.New(
        gateway.WithConfig(cfg),
        gateway.WithLogger(logger),
        gateway.WithAuthService(authService),
        gateway.WithRateLimiter(rateLimiter),
    )

    // Setup routes
    router := gin.New()
    router.Use(gin.Recovery())
    router.Use(middleware.Logger(logger))
    router.Use(middleware.CORS())
    
    gw.SetupRoutes(router)

    // Create server
    srv := &http.Server{
        Addr:         cfg.Server.Address,
        Handler:      router,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    // Start server
    go func() {
        logger.Info("Starting server", zap.String("address", srv.Addr))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Failed to start server", zap.Error(err))
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("Shutting down server...")

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Error("Server forced to shutdown", zap.Error(err))
    }

    logger.Info("Server exited")
}