package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthHandler checks database and cache connection status dynamically and returns diagnostic JSON.
// @Summary Health Check
// @Description Checks PostgreSQL and Redis connection statuses and returns overall health status.
// @Tags General
// @Produce json
// @Success 200 {object} map[string]string "healthy status response"
// @Failure 500 {object} map[string]string "unhealthy status response"
// @Router /healthz [get]
func HealthHandler(dbConn *gorm.DB, rdbClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbStatus := "up"
		redisStatus := "up"
		hasError := false

		// 1. Ping PostgreSQL database
		sqlDB, err := dbConn.DB()
		if err != nil {
			slog.Error("Failed to get sql.DB from gorm", slog.String("error", err.Error()))
			dbStatus = "down"
			hasError = true
		} else {
			// ping database with 3 second timeout context
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			if err := sqlDB.PingContext(ctx); err != nil {
				slog.Error("PostgreSQL ping failure", slog.String("error", err.Error()))
				dbStatus = "down"
				hasError = true
			}
		}

		// 2. Ping Redis cache
		// ping Redis with 3 second timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := rdbClient.Ping(ctx).Err(); err != nil {
			slog.Error("Redis ping failure", slog.String("error", err.Error()))
			redisStatus = "down"
			hasError = true
		}

		// 3. Return response code based on check results
		status := "healthy"
		httpStatus := 200
		if hasError {
			status = "unhealthy"
			httpStatus = 500
		}

		c.JSON(httpStatus, gin.H{
			"status":   status,
			"database": dbStatus,
			"cache":    redisStatus,
		})
	}
}
