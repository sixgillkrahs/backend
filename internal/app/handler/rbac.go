package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/backend/internal/app/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- Roles Handlers ---

func GetRoles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var roles []model.Role
		if err := db.Find(&roles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, roles)
	}
}

func CreateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.Role
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	}
}

func UpdateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}
		var role model.Role
		if err := db.First(&role, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		var req model.Role
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		role.Name = req.Name
		role.Description = req.Description
		if err := db.Save(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, role)
	}
}

func DeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}
		if err := db.Delete(&model.Role{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
	}
}

// --- Resources Handlers ---

func GetResources(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resources []model.Resource
		if err := db.Find(&resources).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resources)
	}
}

func CreateResource(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.Resource
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	}
}

func UpdateResource(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
			return
		}
		var res model.Resource
		if err := db.First(&res, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
			return
		}
		var req model.Resource
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res.Name = req.Name
		res.Description = req.Description
		if err := db.Save(&res).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func DeleteResource(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
			return
		}
		if err := db.Delete(&model.Resource{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
	}
}

// --- Actions Handlers ---

func GetActions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var actions []model.Action
		if err := db.Find(&actions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, actions)
	}
}

func CreateAction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.Action
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	}
}

func UpdateAction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
			return
		}
		var act model.Action
		if err := db.First(&act, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
			return
		}
		var req model.Action
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		act.Name = req.Name
		act.Description = req.Description
		if err := db.Save(&act).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, act)
	}
}

func DeleteAction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
			return
		}
		if err := db.Delete(&model.Action{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Action deleted successfully"})
	}
}

// --- Policies Handlers ---

func GetPolicies(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var policies []model.Policy
		// Preload Role, Resource, and Action definitions for rich output representation
		err := db.Preload("Role").Preload("Resource").Preload("Action").Find(&policies).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, policies)
	}
}

func CreatePolicy(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.Policy
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload relationship fields for response
		db.Preload("Role").Preload("Resource").Preload("Action").First(&req, req.ID)
		c.JSON(http.StatusCreated, req)
	}
}

func UpdatePolicy(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID"})
			return
		}
		var policy model.Policy
		if err := db.First(&policy, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
			return
		}
		var req model.Policy
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		policy.RoleID = req.RoleID
		policy.ResourceID = req.ResourceID
		policy.ActionID = req.ActionID
		policy.Effect = req.Effect
		if err := db.Save(&policy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		db.Preload("Role").Preload("Resource").Preload("Action").First(&policy, policy.ID)
		c.JSON(http.StatusOK, policy)
	}
}

func DeletePolicy(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID"})
			return
		}
		if err := db.Delete(&model.Policy{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Policy deleted successfully"})
	}
}

// --- User-Role Assignments ---

type UserRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

func AssignUserRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req UserRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify role and user exist
		var role model.Role
		if err := db.First(&role, req.RoleID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		userRole := model.UserRole{
			UserID: uint(userID),
			RoleID: req.RoleID,
		}

		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole).Error; err != nil {
			slog.Error("Failed to map user to role", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role to user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
	}
}

func RemoveUserRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req UserRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = db.Where("user_id = ? AND role_id = ?", userID, req.RoleID).
			Delete(&model.UserRole{}).Error

		if err != nil {
			slog.Error("Failed to remove user role mapping", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke role from user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Role revoked successfully"})
	}
}
