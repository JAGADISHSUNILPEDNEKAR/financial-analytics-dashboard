package gateway

import (
    "net/http"
    
    "github.com/financial-analytics/api-gateway/internal/config"
    "github.com/financial-analytics/api-gateway/internal/handlers"
    "github.com/financial-analytics/api-gateway/internal/middleware"
    "github.com/financial-analytics/api-gateway/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "go.uber.org/zap"
)

type Gateway struct {
    config       *config.Config
    logger       *zap.Logger
    authService  services.AuthService
    rateLimiter  middleware.RateLimiter
    wsHub        *handlers.WebSocketHub
    upgrader     websocket.Upgrader
}

type Option func(*Gateway)

func New(opts ...Option) *Gateway {
    g := &Gateway{
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                // Configure origin checking for production
                return true
            },
        },
    }
    
    for _, opt := range opts {
        opt(g)
    }
    
    // Initialize WebSocket hub
    g.wsHub = handlers.NewWebSocketHub(g.logger)
    go g.wsHub.Run()
    
    return g
}

func (g *Gateway) SetupRoutes(router *gin.Engine) {
    // Health check
    router.GET("/health", g.handleHealthCheck)
    
    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Public routes
        auth := v1.Group("/auth")
        {
            auth.POST("/login", g.handleLogin)
            auth.POST("/register", g.handleRegister)
            auth.POST("/refresh", g.handleRefreshToken)
        }
        
        // Protected routes
        protected := v1.Group("/")
        protected.Use(middleware.Auth(g.authService))
        protected.Use(middleware.RateLimit(g.rateLimiter))
        {
            // Dashboard routes
            dashboards := protected.Group("/dashboards")
            {
                dashboards.GET("", g.handleGetDashboards)
                dashboards.POST("", g.handleCreateDashboard)
                dashboards.GET("/:id", g.handleGetDashboard)
                dashboards.PUT("/:id", g.handleUpdateDashboard)
                dashboards.DELETE("/:id", g.handleDeleteDashboard)
                dashboards.POST("/:id/share", g.handleShareDashboard)
            }
            
            // Analytics routes
            analytics := protected.Group("/analytics")
            {
                analytics.GET("/indicators/:symbol", g.handleGetIndicators)
                analytics.POST("/calculate", g.handleCalculate)
                analytics.GET("/historical/:symbol", g.handleGetHistorical)
            }
            
            // WebSocket endpoint
            protected.GET("/ws", g.handleWebSocket)
            
            // User routes
            users := protected.Group("/users")
            {
                users.GET("/profile", g.handleGetProfile)
                users.PUT("/profile", g.handleUpdateProfile)
                users.GET("/preferences", g.handleGetPreferences)
                users.PUT("/preferences", g.handleUpdatePreferences)
            }
            
            // Watchlist routes
            watchlists := protected.Group("/watchlists")
            {
                watchlists.GET("", g.handleGetWatchlists)
                watchlists.POST("", g.handleCreateWatchlist)
                watchlists.PUT("/:id", g.handleUpdateWatchlist)
                watchlists.DELETE("/:id", g.handleDeleteWatchlist)
            }
            
            // Alert routes
            alerts := protected.Group("/alerts")
            {
                alerts.GET("", g.handleGetAlerts)
                alerts.POST("", g.handleCreateAlert)
                alerts.PUT("/:id", g.handleUpdateAlert)
                alerts.DELETE("/:id", g.handleDeleteAlert)
            }
        }
    }
}

func (g *Gateway) handleHealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "timestamp": time.Now().Unix(),
    })
}

func (g *Gateway) handleWebSocket(c *gin.Context) {
    conn, err := g.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        g.logger.Error("Failed to upgrade WebSocket", zap.Error(err))
        return
    }
    
    userID := c.GetString("user_id")
    client := handlers.NewClient(conn, userID, g.wsHub)
    
    g.wsHub.Register <- client
    
    go client.WritePump()
    go client.ReadPump()
}