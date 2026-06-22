package middleware

import (
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
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

type AnyValue struct{}

func (a AnyValue) Match(v driver.Value) bool {
	return true
}

func TestRBACMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setupRouter := func(db *gorm.DB, resource, action string, setupContext func(c *gin.Context)) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			if setupContext != nil {
				setupContext(c)
			}
			c.Next()
		})
		r.Use(RBACMiddleware(db, resource, action))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})
		return r
	}

	t.Run("Missing User ID Context", func(t *testing.T) {
		db, _, err := NewMockDB()
		assert.NoError(t, err)

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			// Do not set userID
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authentication required")
	})

	t.Run("DB Error Fetching User Roles", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(1)
		mock.ExpectQuery(`SELECT "role_id" FROM "user_roles" WHERE user_id = .*`).
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			c.Set("userID", userID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal database error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No Roles Assigned", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(1)
		mock.ExpectQuery(`SELECT "role_id" FROM "user_roles" WHERE user_id = .*`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"role_id"}))

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			c.Set("userID", userID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Access denied: no roles assigned to user")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Resource Not Found", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(1)
		mock.ExpectQuery(`SELECT "role_id" FROM "user_roles" WHERE user_id = .*`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"role_id"}).AddRow(10))

		mock.ExpectQuery(`SELECT "id" FROM "resources" WHERE name = .*`).
			WithArgs("reports").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			c.Set("userID", userID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Access denied: requested resource does not exist in configuration")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Action Not Found", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(1)
		mock.ExpectQuery(`SELECT "role_id" FROM "user_roles" WHERE user_id = .*`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"role_id"}).AddRow(10))

		mock.ExpectQuery(`SELECT "id" FROM "resources" WHERE name = .*`).
			WithArgs("reports").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		mock.ExpectQuery(`SELECT "id" FROM "actions" WHERE name = .*`).
			WithArgs("read").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			c.Set("userID", userID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Access denied: requested action does not exist in configuration")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Policy Deny", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(1)
		mock.ExpectQuery(`SELECT "role_id" FROM "user_roles" WHERE user_id = .*`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"role_id"}).AddRow(10))

		mock.ExpectQuery(`SELECT "id" FROM "resources" WHERE name = .*`).
			WithArgs("reports").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		mock.ExpectQuery(`SELECT "id" FROM "actions" WHERE name = .*`).
			WithArgs("read").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

		mock.ExpectQuery(`SELECT effect FROM "policies" WHERE role_id IN .*`).
			WithArgs(10, 2, 3).
			WillReturnRows(sqlmock.NewRows([]string{"effect"}).AddRow("deny"))

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			c.Set("userID", userID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Access denied: explicitly forbidden by security policy")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Policy Allow", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(1)
		mock.ExpectQuery(`SELECT "role_id" FROM "user_roles" WHERE user_id = .*`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"role_id"}).AddRow(10).AddRow(11))

		mock.ExpectQuery(`SELECT "id" FROM "resources" WHERE name = .*`).
			WithArgs("reports").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		mock.ExpectQuery(`SELECT "id" FROM "actions" WHERE name = .*`).
			WithArgs("read").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

		mock.ExpectQuery(`SELECT effect FROM "policies" WHERE role_id IN .*`).
			WithArgs(10, 11, 2, 3).
			WillReturnRows(sqlmock.NewRows([]string{"effect"}).AddRow("allow"))

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			c.Set("userID", userID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"success"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Policy Default Deny", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		userID := uint(1)
		mock.ExpectQuery(`SELECT "role_id" FROM "user_roles" WHERE user_id = .*`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"role_id"}).AddRow(10))

		mock.ExpectQuery(`SELECT "id" FROM "resources" WHERE name = .*`).
			WithArgs("reports").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		mock.ExpectQuery(`SELECT "id" FROM "actions" WHERE name = .*`).
			WithArgs("read").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

		mock.ExpectQuery(`SELECT effect FROM "policies" WHERE role_id IN .*`).
			WithArgs(10, 2, 3).
			WillReturnRows(sqlmock.NewRows([]string{"effect"}))

		r := setupRouter(db, "reports", "read", func(c *gin.Context) {
			c.Set("userID", userID)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Access denied: no matching allow policies found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
