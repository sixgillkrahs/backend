package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/username/backend/internal/app/model"
	"github.com/username/backend/internal/pkg/auth"
	"github.com/username/backend/internal/pkg/config"
	"gorm.io/gorm"
)

// SignupHandler registers a new user with email, name, and hashed password.
func SignupHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if email already exists
		var existingUser model.User
		err := db.Where("email = ?", req.Email).First(&existingUser).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
			return
		} else if err != gorm.ErrRecordNotFound {
			slog.Error("Database query error on signup", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Password hashing failure", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		user := model.User{
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: string(passwordHash),
		}

		if err := db.Create(&user).Error; err != nil {
			slog.Error("Failed to create user in DB", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"user":    user,
		})
	}
}

// LoginHandler authenticates credentials and returns an IP-locked JWT.
func LoginHandler(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User
		err := db.Where("email = ?", req.Email).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
				return
			}
			slog.Error("Database query error on login", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate token with current client IP
		clientIP := c.ClientIP()
		token, err := auth.GenerateToken(user.ID, clientIP, cfg.JWTSecret, cfg.JWTExpirationHours)
		if err != nil {
			slog.Error("JWT generation failed", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, model.LoginResponse{
			Token: token,
			User:  user,
		})
	}
}
