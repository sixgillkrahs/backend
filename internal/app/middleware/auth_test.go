package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sixgillkrahs/backend/internal/pkg/auth"
	"github.com/sixgillkrahs/backend/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret:          "my-jwt-test-secret",
		JWTExpirationHours: 2,
	}

	// Helper to setup mock router
	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(AuthMiddleware(cfg))
		r.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			clientIP, _ := c.Get("clientIP")
			c.JSON(http.StatusOK, gin.H{
				"userID":   userID,
				"clientIP": clientIP,
			})
		})
		return r
	}

	// 1. Missing header
	t.Run("Missing Header", func(t *testing.T) {
		r := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	// 2. Invalid format (not Bearer)
	t.Run("Invalid Format", func(t *testing.T) {
		r := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Basic credentials")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header must be Bearer token")
	})

	// 3. Invalid token
	t.Run("Invalid Token", func(t *testing.T) {
		r := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalidtokenhere")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})

	// 4. IP Mismatch
	t.Run("IP Mismatch", func(t *testing.T) {
		token, err := auth.GenerateToken(10, "192.168.1.1", cfg.JWTSecret, cfg.JWTExpirationHours)
		assert.NoError(t, err)

		r := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.RemoteAddr = "192.168.1.2:1234" // causes mismatch
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Session invalidated: client IP changed")
	})

	// 5. Success
	t.Run("Success", func(t *testing.T) {
		clientIP := "192.0.2.1"
		token, err := auth.GenerateToken(42, clientIP, cfg.JWTSecret, cfg.JWTExpirationHours)
		assert.NoError(t, err)

		r := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.RemoteAddr = clientIP + ":1234"
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"userID":42`)
		assert.Contains(t, w.Body.String(), `"clientIP":"192.0.2.1"`)
	})

	// 6. Success with Query Parameter Token
	t.Run("Success with Query Parameter Token", func(t *testing.T) {
		clientIP := "192.0.2.1"
		token, err := auth.GenerateToken(100, clientIP, cfg.JWTSecret, cfg.JWTExpirationHours)
		assert.NoError(t, err)

		r := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test?token="+token, nil)
		req.RemoteAddr = clientIP + ":1234"
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"userID":100`)
		assert.Contains(t, w.Body.String(), `"clientIP":"192.0.2.1"`)
	})
}
