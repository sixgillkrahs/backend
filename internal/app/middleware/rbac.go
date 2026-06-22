package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RBACMiddleware checks if the authenticated user has permission to access a resource and action.
// It queries the roles, resources, actions, and policies tables dynamically.
func RBACMiddleware(db *gorm.DB, resourceName string, actionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve userID set by AuthMiddleware
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userID, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID context format"})
			c.Abort()
			return
		}

		// 1. Query the user's role IDs
		var roleIDs []uint
		err := db.Table("user_roles").
			Where("user_id = ?", userID).
			Pluck("role_id", &roleIDs).Error

		if err != nil {
			slog.Error("Failed to fetch user roles from DB", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal database error"})
			c.Abort()
			return
		}

		if len(roleIDs) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: no roles assigned to user"})
			c.Abort()
			return
		}

		// 2. Resolve resource ID
		var resourceID uint
		err = db.Table("resources").
			Where("name = ?", resourceName).
			Pluck("id", &resourceID).Error

		if err != nil {
			slog.Error("Failed to query resource ID", slog.String("resource", resourceName), slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal database error"})
			c.Abort()
			return
		}
		if resourceID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: requested resource does not exist in configuration"})
			c.Abort()
			return
		}

		// 3. Resolve action ID
		var actionID uint
		err = db.Table("actions").
			Where("name = ?", actionName).
			Pluck("id", &actionID).Error

		if err != nil {
			slog.Error("Failed to query action ID", slog.String("action", actionName), slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal database error"})
			c.Abort()
			return
		}
		if actionID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: requested action does not exist in configuration"})
			c.Abort()
			return
		}

		// 4. Query policies matching the user's roles, resource, and action
		type policyResult struct {
			Effect string
		}
		var policies []policyResult
		err = db.Table("policies").
			Select("effect").
			Where("role_id IN (?) AND resource_id = ? AND action_id = ?", roleIDs, resourceID, actionID).
			Scan(&policies).Error

		if err != nil {
			slog.Error("Failed to query policy checks", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal database error"})
			c.Abort()
			return
		}

		// Evaluate policy decisions:
		// - If there is ANY 'deny' rule, access is blocked immediately.
		// - If there is at least one 'allow' rule and no 'deny' rules, access is granted.
		// - Otherwise, access is blocked (default deny).
		allowed := false
		for _, p := range policies {
			if p.Effect == "deny" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: explicitly forbidden by security policy"})
				c.Abort()
				return
			}
			if p.Effect == "allow" {
				allowed = true
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: no matching allow policies found"})
			c.Abort()
			return
		}

		c.Next()
	}
}
