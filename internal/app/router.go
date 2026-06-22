package app

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/username/backend/internal/app/handler"
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

	// Register GET /ping route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Register GET /healthz route
	r.GET("/healthz", handler.HealthHandler(dbConn, rdbClient))

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
