package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/username/backend/internal/app/model"
	"gorm.io/gorm"
)

func TestGetRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Get Roles", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		mock.ExpectQuery(`SELECT .* FROM "roles"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "admin", "Admin role").
				AddRow(2, "member", "Member role"))

		r := gin.New()
		r.GET("/roles", GetRoles(db))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"admin"`)
		assert.Contains(t, w.Body.String(), `"name":"member"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DB Error Get Roles", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		mock.ExpectQuery(`SELECT .* FROM "roles"`).
			WillReturnError(errors.New("roles query failed"))

		r := gin.New()
		r.GET("/roles", GetRoles(db))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "roles query failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Create Role", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		newRole := model.Role{
			Name:        "editor",
			Description: "Editor role",
		}
		bodyBytes, _ := json.Marshal(newRole)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "roles" .*`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
		mock.ExpectCommit()

		r := gin.New()
		r.POST("/roles", CreateRole(db))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"editor"`)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAssignUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Assign Role", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		roleID := uint(10)
		userID := 5
		reqBody := UserRoleRequest{
			RoleID: roleID,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		// 1. First role check
		mock.ExpectQuery(`SELECT .* FROM "roles" WHERE .*id.* = .* LIMIT 1?`).
			WithArgs(roleID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(roleID, "admin"))

		// 2. First user check
		mock.ExpectQuery(`SELECT .* FROM "users" WHERE .*id.* = .* LIMIT 1?`).
			WithArgs(userID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(userID, "Bob"))

		// 3. Insert user_roles mapping
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "user_roles" .*`).
			WithArgs(userID, roleID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		r := gin.New()
		r.POST("/users/:id/roles", AssignUserRole(db))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/5/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Role assigned successfully")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Role Not Found", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		roleID := uint(99)
		reqBody := UserRoleRequest{
			RoleID: roleID,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mock.ExpectQuery(`SELECT .* FROM "roles" WHERE .*id.* = .* LIMIT 1?`).
			WithArgs(roleID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		r := gin.New()
		r.POST("/users/:id/roles", AssignUserRole(db))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/5/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Role not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRemoveUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Remove Role", func(t *testing.T) {
		db, mock, err := NewMockDB()
		assert.NoError(t, err)

		roleID := uint(10)
		userID := 5
		reqBody := UserRoleRequest{
			RoleID: roleID,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "user_roles" WHERE .*user_id = .* AND role_id = .*`).
			WithArgs(userID, roleID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		r := gin.New()
		r.DELETE("/users/:id/roles", RemoveUserRole(db))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/5/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Role revoked successfully")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
