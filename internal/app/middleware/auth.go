package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/username/backend/internal/pkg/auth"
	"github.com/username/backend/internal/pkg/config"
)

// AuthMiddleware extracts the JWT token from the Authorization header and validates it.
// It also enforces IP locking by checking if the client IP has changed.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer token"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			slog.Warn("JWT validation failed", slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// IP Lock check
		currentIP := c.ClientIP()
		if claims.IP != currentIP {
			slog.Warn("JWT IP mismatch",
				slog.Uint64("user_id", uint64(claims.UserID)),
				slog.String("token_ip", claims.IP),
				slog.String("request_ip", currentIP),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalidated: client IP changed"})
			c.Abort()
			return
		}

		// Store user info in context for downstream handlers/middlewares
		c.Set("userID", claims.UserID)
		c.Set("clientIP", claims.IP)

		c.Next()
	}
}
