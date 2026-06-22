package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func NewMockDBWithPing() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		return nil, nil, err
	}
	mock.ExpectPing() // Consume GORM's startup ping
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Health Status", func(t *testing.T) {
		// Mock DB
		gormDB, mock, err := NewMockDBWithPing()
		assert.NoError(t, err)
		mock.ExpectPing()

		// Mock Redis
		mr, err := miniredis.Run()
		assert.NoError(t, err)
		defer mr.Close()

		rdbClient := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
		defer rdbClient.Close()

		// Setup Gin router
		r := gin.New()
		r.GET("/healthz", HealthHandler(gormDB, rdbClient))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"healthy"`)
		assert.Contains(t, w.Body.String(), `"database":"up"`)
		assert.Contains(t, w.Body.String(), `"cache":"up"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PostgreSQL Down Health Status", func(t *testing.T) {
		// Mock DB returning ping error
		gormDB, mock, err := NewMockDBWithPing()
		assert.NoError(t, err)
		mock.ExpectPing().WillReturnError(errors.New("postgres connection failed"))

		// Mock Redis
		mr, err := miniredis.Run()
		assert.NoError(t, err)
		defer mr.Close()

		rdbClient := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
		defer rdbClient.Close()

		r := gin.New()
		r.GET("/healthz", HealthHandler(gormDB, rdbClient))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"unhealthy"`)
		assert.Contains(t, w.Body.String(), `"database":"down"`)
		assert.Contains(t, w.Body.String(), `"cache":"up"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Redis Down Health Status", func(t *testing.T) {
		// Mock DB
		gormDB, mock, err := NewMockDBWithPing()
		assert.NoError(t, err)
		mock.ExpectPing()

		// Redis client pointing to invalid address
		rdbClient := redis.NewClient(&redis.Options{
			Addr: "localhost:9999", // closed port
		})
		defer rdbClient.Close()

		r := gin.New()
		r.GET("/healthz", HealthHandler(gormDB, rdbClient))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"unhealthy"`)
		assert.Contains(t, w.Body.String(), `"database":"up"`)
		assert.Contains(t, w.Body.String(), `"cache":"down"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
