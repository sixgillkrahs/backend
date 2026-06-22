package main

import (
	"log/slog"
	"os"

	"github.com/username/backend/internal/pkg/cache"
	"github.com/username/backend/internal/pkg/config"
	"github.com/username/backend/internal/pkg/db"
)

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

	_ = database // Will be used by handlers/repositories in future phases

	// Initialize Redis Cache Client
	redisClient, err := cache.InitRedis(cfg)
	if err != nil {
		slog.Error("Redis initialization failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	_ = redisClient // Will be used by handlers/repositories/caches in future phases
}
