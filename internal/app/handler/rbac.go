package handler

import (
	"errors"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sixgillkrahs/backend/internal/app/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- Roles Handlers ---

// GetRoles retrieves all configured roles from the database.
// @Summary Get Roles
// @Description Retrieves all roles configured in the system.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Role "List of roles"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/roles [get]
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

// CreateRole creates a new dynamic role in the database.
// @Summary Create Role
// @Description Creates a new custom role.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.Role true "Role creation data"
// @Success 201 {object} model.Role "Role created successfully"
// @Failure 400 {object} map[string]string "Bad request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/roles [post]
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

// UpdateRole updates an existing role configuration by ID.
// @Summary Update Role
// @Description Updates the name or description of a role.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body model.Role true "Role update data"
// @Success 200 {object} model.Role "Role updated successfully"
// @Failure 400 {object} map[string]string "Invalid ID or bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/roles/{id} [put]
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

// DeleteRole deletes a role by ID.
// @Summary Delete Role
// @Description Deletes a role and revokes its policy mappings.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string "Role deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/roles/{id} [delete]
func DeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}

		var role model.Role
		if err := db.First(&role, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// Prevent deleting system default roles
		if role.Name == "admin" || role.Name == "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete system default role: " + role.Name})
			return
		}

		if err := db.Delete(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
	}
}

// --- Resources Handlers ---

// GetResources retrieves all system resource definitions.
// @Summary Get Resources
// @Description Retrieves all resources defined in the system.
// @Tags Resources
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Resource "List of resources"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/resources [get]
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

// CreateResource creates a new system resource definition.
// @Summary Create Resource
// @Description Creates a new resource descriptor.
// @Tags Resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.Resource true "Resource creation data"
// @Success 201 {object} model.Resource "Resource created successfully"
// @Failure 400 {object} map[string]string "Bad request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/resources [post]
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

// UpdateResource updates an existing resource configuration by ID.
// @Summary Update Resource
// @Description Updates the name or description of a resource.
// @Tags Resources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Resource ID"
// @Param request body model.Resource true "Resource update data"
// @Success 200 {object} model.Resource "Resource updated successfully"
// @Failure 400 {object} map[string]string "Invalid ID or bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 404 {object} map[string]string "Resource not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/resources/{id} [put]
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

// DeleteResource deletes a resource definition by ID.
// @Summary Delete Resource
// @Description Deletes a resource definition.
// @Tags Resources
// @Produce json
// @Security BearerAuth
// @Param id path int true "Resource ID"
// @Success 200 {object} map[string]string "Resource deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/resources/{id} [delete]
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

// GetActions retrieves all authorization actions.
// @Summary Get Actions
// @Description Retrieves all action descriptors defined in the system.
// @Tags Actions
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Action "List of actions"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/actions [get]
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

// CreateAction creates a new dynamic authorization action.
// @Summary Create Action
// @Description Creates a new action descriptor.
// @Tags Actions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.Action true "Action creation data"
// @Success 201 {object} model.Action "Action created successfully"
// @Failure 400 {object} map[string]string "Bad request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/actions [post]
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

// UpdateAction updates an existing action configuration by ID.
// @Summary Update Action
// @Description Updates the name or description of an action.
// @Tags Actions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Action ID"
// @Param request body model.Action true "Action update data"
// @Success 200 {object} model.Action "Action updated successfully"
// @Failure 400 {object} map[string]string "Invalid ID or bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 404 {object} map[string]string "Action not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/actions/{id} [put]
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

// DeleteAction deletes an action definition by ID.
// @Summary Delete Action
// @Description Deletes an action definition.
// @Tags Actions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Action ID"
// @Success 200 {object} map[string]string "Action deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/actions/{id} [delete]
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

// GetPolicies retrieves all active access policies.
// @Summary Get Policies
// @Description Retrieves all policies. Preloads Roles, Resources, and Actions relationship structs.
// @Tags Policies
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Policy "List of policies"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/policies [get]
func GetPolicies(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var policies []model.Policy
		err := db.Preload("Role").Preload("Resource").Preload("Action").Find(&policies).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, policies)
	}
}

// CreatePolicy creates a new dynamic access control policy.
// @Summary Create Policy
// @Description Creates a new policy mapping a role to a resource and action.
// @Tags Policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.Policy true "Policy mapping data"
// @Success 201 {object} model.Policy "Policy created successfully"
// @Failure 400 {object} map[string]string "Bad request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/policies [post]
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
		db.Preload("Role").Preload("Resource").Preload("Action").First(&req, req.ID)
		c.JSON(http.StatusCreated, req)
	}
}

// UpdatePolicy updates an existing policy configuration by ID.
// @Summary Update Policy
// @Description Updates the role, resource, action, or effect of a policy.
// @Tags Policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Policy ID"
// @Param request body model.Policy true "Policy update data"
// @Success 200 {object} model.Policy "Policy updated successfully"
// @Failure 400 {object} map[string]string "Invalid ID or bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 404 {object} map[string]string "Policy not found"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/policies/{id} [put]
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

// DeletePolicy deletes a policy by ID.
// @Summary Delete Policy
// @Description Deletes a policy definition.
// @Tags Policies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Policy ID"
// @Success 200 {object} map[string]string "Policy deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/policies/{id} [delete]
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

// AssignUserRole assigns a role to a user.
// @Summary Assign User Role
// @Description Maps a user to a dynamic role.
// @Tags User Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UserRoleRequest true "Role mapping data"
// @Success 200 {object} map[string]string "Role assigned successfully"
// @Failure 400 {object} map[string]string "Invalid ID or request payload"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 404 {object} map[string]string "User or Role not found"
// @Failure 500 {object} map[string]string "Internal database error"
// @Router /api/v1/users/{id}/roles [post]
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

// RemoveUserRole revokes a role from a user.
// @Summary Revoke User Role
// @Description Removes a user role mapping.
// @Tags User Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UserRoleRequest true "Role revocation data"
// @Success 200 {object} map[string]string "Role revoked successfully"
// @Failure 400 {object} map[string]string "Invalid ID or request payload"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Internal database error"
// @Router /api/v1/users/{id}/roles [delete]
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

// GetUsers returns a paginated, sorted, and filtered list of users with their primary role.
// @Summary List Users
// @Description Returns a list of users supporting search, role/status filtering, pagination, and sorting.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort_by query string false "Sort by field (name, email, status, role, joinedDate)" default(name)
// @Param sort_order query string false "Sort order (asc, desc)" default(asc)
// @Param search query string false "Search query filtering name or email"
// @Param role query string false "Filter by role name"
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{} "Paginated list of users"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden by RBAC policy"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/users [get]
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		sortBy := c.DefaultQuery("sort_by", "name")
		sortOrder := c.DefaultQuery("sort_order", "asc")
		search := c.Query("search")
		roleFilter := c.Query("role")
		statusFilter := c.Query("status")

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}

		offset := (page - 1) * limit

		query := db.Table("users").
			Select("users.id, users.name, users.email, users.status, COALESCE(roles.name, 'User') as role, users.created_at").
			Joins("LEFT JOIN user_roles ON user_roles.user_id = users.id").
			Joins("LEFT JOIN roles ON roles.id = user_roles.role_id")

		if search != "" {
			query = query.Where("users.name ILIKE ? OR users.email ILIKE ?", "%"+search+"%", "%"+search+"%")
		}
		if roleFilter != "" {
			if roleFilter == "User" {
				query = query.Where("roles.name = ? OR roles.name IS NULL", roleFilter)
			} else {
				query = query.Where("roles.name = ?", roleFilter)
			}
		}
		if statusFilter != "" {
			query = query.Where("users.status = ?", statusFilter)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		orderClause := "users.name"
		if sortBy == "email" {
			orderClause = "users.email"
		} else if sortBy == "status" {
			orderClause = "users.status"
		} else if sortBy == "role" {
			orderClause = "roles.name"
		} else if sortBy == "joinedDate" {
			orderClause = "users.created_at"
		}

		if sortOrder == "desc" {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}

		type UserRow struct {
			ID        uint      `json:"id"`
			Name      string    `json:"name"`
			Email     string    `json:"email"`
			Role      string    `json:"role"`
			Status    string    `json:"status"`
			CreatedAt time.Time `json:"created_at"`
		}

		var rows []UserRow
		if err := query.Order(orderClause).Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type UserResponse struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Email      string `json:"email"`
			Role       string `json:"role"`
			Status     string `json:"status"`
			JoinedDate string `json:"joinedDate"`
		}

		usersList := make([]UserResponse, len(rows))
		for i, r := range rows {
			usersList[i] = UserResponse{
				ID:         strconv.FormatUint(uint64(r.ID), 10),
				Name:       r.Name,
				Email:      r.Email,
				Role:       r.Role,
				Status:     r.Status,
				JoinedDate: r.CreatedAt.Format("2006-01-02"),
			}
		}

		pages := int(math.Ceil(float64(total) / float64(limit)))

		c.JSON(http.StatusOK, gin.H{
			"users": usersList,
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": pages,
		})
	}
}
