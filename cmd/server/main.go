package main

import (
	"log/slog"
	"os"

	"github.com/username/backend/internal/pkg/config"
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
}
