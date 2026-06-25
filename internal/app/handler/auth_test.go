package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sixgillkrahs/backend/internal/app/model"
	"github.com/sixgillkrahs/backend/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestSignupHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Signup", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		reqBody := model.SignupRequest{
			Name:     "Alice Smith",
			Email:    "alice@example.com",
			Password: "password123",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		// Email check: expects user not found
		mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = .* LIMIT 1?`).
			WithArgs(reqBody.Email, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Creation: expects insert statement
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users" .*`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		r := gin.New()
		r.POST("/signup", NewAuthHandler(db, config.LoadConfig()).SignupHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered successfully")
		assert.Contains(t, w.Body.String(), `"email":"alice@example.com"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Email Conflict", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		reqBody := model.SignupRequest{
			Name:     "Alice Smith",
			Email:    "alice@example.com",
			Password: "password123",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		// Email check finds existing user
		mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = .* LIMIT 1?`).
			WithArgs(reqBody.Email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, reqBody.Email))

		r := gin.New()
		r.POST("/signup", NewAuthHandler(db, config.LoadConfig()).SignupHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "Email is already registered")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Bad Request JSON Payload", func(t *testing.T) {
		db, _, err := NewMockDB()
		assert.NoError(t, err)

		r := gin.New()
		r.POST("/signup", NewAuthHandler(db, config.LoadConfig()).SignupHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(`{"name":""}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret:          "secret-key",
		JWTExpirationHours: 2,
	}

	t.Run("Success Login", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		reqBody := model.LoginRequest{
			Email:    "alice@example.com",
			Password: password,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		// Email check finds user record
		mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = .* LIMIT 1?`).
			WithArgs(reqBody.Email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
				AddRow(42, reqBody.Email, string(hashedPassword)))

		r := gin.New()
		r.POST("/login", NewAuthHandler(db, cfg).LoginHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"token":`)
		assert.Contains(t, w.Body.String(), `"email":"alice@example.com"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User Not Found", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		reqBody := model.LoginRequest{
			Email:    "unknown@example.com",
			Password: "password123",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = .* LIMIT 1?`).
			WithArgs(reqBody.Email, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		r := gin.New()
		r.POST("/login", NewAuthHandler(db, config.LoadConfig()).LoginHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		correctPassword := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

		reqBody := model.LoginRequest{
			Email:    "alice@example.com",
			Password: "wrongpassword",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = .* LIMIT 1?`).
			WithArgs(reqBody.Email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
				AddRow(42, reqBody.Email, string(hashedPassword)))

		r := gin.New()
		r.POST("/login", NewAuthHandler(db, cfg).LoginHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Get Profile", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(42)
		clientIP := "192.168.1.50"

		mock.ExpectQuery(`SELECT .* FROM "users" WHERE .*id.* = .* LIMIT 1?`).
			WithArgs(userID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(userID, "Alice Smith", "alice@example.com"))

		r := gin.New()
		r.GET("/profile", func(c *gin.Context) {
			c.Set("userID", userID)
			c.Set("clientIP", clientIP)
			c.Next()
		}, NewAuthHandler(db, config.LoadConfig()).ProfileHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Alice Smith"`)
		assert.Contains(t, w.Body.String(), `"token_ip":"192.168.1.50"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User Missing DB Record", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(99)

		mock.ExpectQuery(`SELECT .* FROM "users" WHERE .*id.* = .* LIMIT 1?`).
			WithArgs(userID, 1).
			WillReturnError(errors.New("record not found"))

		r := gin.New()
		r.GET("/profile", func(c *gin.Context) {
			c.Set("userID", userID)
			c.Next()
		}, NewAuthHandler(db, config.LoadConfig()).ProfileHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Authentication Context Missing", func(t *testing.T) {
		db, _, err := NewMockDB()
		assert.NoError(t, err)

		r := gin.New()
		r.GET("/profile", NewAuthHandler(db, config.LoadConfig()).ProfileHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authentication required")
	})
}
