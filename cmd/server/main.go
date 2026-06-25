package main

import (
	"log/slog"
	"os"

	"github.com/sixgillkrahs/backend/internal/app"
	"github.com/sixgillkrahs/backend/internal/pkg/cache"
	"github.com/sixgillkrahs/backend/internal/pkg/config"
	"github.com/sixgillkrahs/backend/internal/pkg/db"
)

// @title Golang Gin Backend Scaffold API
// @version 1.0
// @description Interactive API documentation with dynamic RBAC and IP-locked JWT.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your-token>" to authenticate.
func main() {
	// Initialize structured logging using slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.LoadConfig()

	slog.Info("Starting server",
		slog.String("port", cfg.Port),
		slog.String("env", cfg.Env),
	)

	// Initialize PostgreSQL Connection Pool and Run Migrations
	database, err := db.InitPostgres(cfg)
	if err != nil {
		slog.Error("Database initialization failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize Redis Cache Client
	redisClient, err := cache.InitRedis(cfg)
	if err != nil {
		slog.Error("Redis initialization failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Setup Gin HTTP Router
	r := app.SetupRouter(cfg, database, redisClient)

	// Run HTTP Server
	slog.Info("HTTP Server is listening", slog.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("Server failed to run", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
