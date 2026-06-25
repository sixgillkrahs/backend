package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sixgillkrahs/backend/internal/app/model"
	"github.com/sixgillkrahs/backend/internal/pkg/auth"
	"github.com/sixgillkrahs/backend/internal/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:  db,
		cfg: cfg,
	}
}

// SignupHandler registers a new user with email, name, and hashed password.
// @Summary User Signup
// @Description Register a new user with their name, email, and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.SignupRequest true "Signup payload"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]string "Bad request validation payload"
// @Failure 409 {object} map[string]string "Email already registered conflict"
// @Failure 500 {object} map[string]string "Internal database or server error"
// @Router /api/v1/auth/signup [post]
func (a *AuthHandler) SignupHandler(c *gin.Context) {
	var req model.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser model.User
	err := a.db.Where("email = ?", req.Email).First(&existingUser).Error
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
		Status:       "active",
	}

	if err := a.db.Create(&user).Error; err != nil {
		slog.Error("Failed to create user in DB", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// LoginHandler authenticates credentials and returns an IP-locked JWT.
// @Summary User Login
// @Description Authenticates email and password, returning an IP-locked JWT access token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials payload"
// @Success 200 {object} model.LoginResponse "Successfully logged in and generated JWT"
// @Failure 400 {object} map[string]string "Bad request validation payload"
// @Failure 401 {object} map[string]string "Invalid email or password credentials"
// @Failure 500 {object} map[string]string "Internal JWT or server error"
// @Router /api/v1/auth/login [post]
func (a *AuthHandler) LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err := a.db.Where("email = ?", req.Email).First(&user).Error
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
	token, err := auth.GenerateToken(user.ID, clientIP, a.cfg.JWTSecret, a.cfg.JWTExpirationHours)
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

// ProfileHandler retrieves the authenticated user's profile info and request/token IP properties.
// @Summary Get User Profile
// @Description Fetches the current user profile from PostgreSQL using the authenticated context. Checks client IP lock.
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Successful profile response"
// @Failure 401 {object} map[string]string "Unauthorized or IP lock validation failure"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/v1/profile [get]
func (a *AuthHandler) ProfileHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	clientIP, _ := c.Get("clientIP")

	var user model.User
	if err := a.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":      user,
		"token_ip":  clientIP,
		"client_ip": c.ClientIP(),
	})
}
