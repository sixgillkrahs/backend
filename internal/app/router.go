package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/username/backend/docs"
	"github.com/username/backend/internal/app/handler"
	"github.com/username/backend/internal/app/middleware"
	"github.com/username/backend/internal/pkg/config"
	"gorm.io/gorm"
)

// SetupRouter initializes the Gin engine with standard recovery and custom slog request logging.
func SetupRouter(cfg *config.Config, dbConn *gorm.DB, rdbClient *redis.Client) *gin.Engine {
	// Set Gin mode dynamically
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Register Recovery middleware
	r.Use(gin.Recovery())

	// Register custom slog logger middleware
	r.Use(slogLogger())

	// Configure CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register GET /ping route
	r.GET("/ping", PingHandler)

	// Register GET /healthz route
	r.GET("/healthz", handler.HealthHandler(dbConn, rdbClient))

	// Register GET /swagger/*any route with automatic redirect from root to index.html
	r.GET("/swagger/*any", func(c *gin.Context) {
		path := c.Param("any")
		if path == "" || path == "/" {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})


	// API v1 Router Group
	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/signup", handler.SignupHandler(dbConn))
			authGroup.POST("/login", handler.LoginHandler(dbConn, cfg))
		}

		// Protected Route Group
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.GET("/profile", handler.ProfileHandler(dbConn))

			// --- RBAC Management (Admin Protected) ---

			// Roles Management
			protected.GET("/roles", middleware.RBACMiddleware(dbConn, "roles", "read"), handler.GetRoles(dbConn))
			protected.POST("/roles", middleware.RBACMiddleware(dbConn, "roles", "create"), handler.CreateRole(dbConn))
			protected.PUT("/roles/:id", middleware.RBACMiddleware(dbConn, "roles", "update"), handler.UpdateRole(dbConn))
			protected.DELETE("/roles/:id", middleware.RBACMiddleware(dbConn, "roles", "delete"), handler.DeleteRole(dbConn))

			// Resources Management
			protected.GET("/resources", middleware.RBACMiddleware(dbConn, "resources", "read"), handler.GetResources(dbConn))
			protected.POST("/resources", middleware.RBACMiddleware(dbConn, "resources", "create"), handler.CreateResource(dbConn))
			protected.PUT("/resources/:id", middleware.RBACMiddleware(dbConn, "resources", "update"), handler.UpdateResource(dbConn))
			protected.DELETE("/resources/:id", middleware.RBACMiddleware(dbConn, "resources", "delete"), handler.DeleteResource(dbConn))

			// Actions Management
			protected.GET("/actions", middleware.RBACMiddleware(dbConn, "actions", "read"), handler.GetActions(dbConn))
			protected.POST("/actions", middleware.RBACMiddleware(dbConn, "actions", "create"), handler.CreateAction(dbConn))
			protected.PUT("/actions/:id", middleware.RBACMiddleware(dbConn, "actions", "update"), handler.UpdateAction(dbConn))
			protected.DELETE("/actions/:id", middleware.RBACMiddleware(dbConn, "actions", "delete"), handler.DeleteAction(dbConn))

			// Policies Management
			protected.GET("/policies", middleware.RBACMiddleware(dbConn, "policies", "read"), handler.GetPolicies(dbConn))
			protected.POST("/policies", middleware.RBACMiddleware(dbConn, "policies", "create"), handler.CreatePolicy(dbConn))
			protected.PUT("/policies/:id", middleware.RBACMiddleware(dbConn, "policies", "update"), handler.UpdatePolicy(dbConn))
			protected.DELETE("/policies/:id", middleware.RBACMiddleware(dbConn, "policies", "delete"), handler.DeletePolicy(dbConn))

			// User Role Mapping
			protected.POST("/users/:id/roles", middleware.RBACMiddleware(dbConn, "users", "update"), handler.AssignUserRole(dbConn))
			protected.DELETE("/users/:id/roles", middleware.RBACMiddleware(dbConn, "users", "update"), handler.RemoveUserRole(dbConn))
		}
	}

	return r
}

func slogLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		slog.Info("Request received",
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		)
	}
}

// PingHandler handles ping pong response.
// @Summary Ping Pong
// @Description returns a simple pong message to verify the server is running.
// @Tags General
// @Produce json
// @Success 200 {object} map[string]string "pong response"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
