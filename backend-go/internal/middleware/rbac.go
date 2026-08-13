package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequirePermission returns a Gin middleware that checks whether the
// authenticated user (via JWT claims) has the specified permission.
//
// Permissions follow the format "resource:action", e.g. "employee:write".
// The super_admin role bypasses all permission checks.
//
// Usage:
//
//	router.POST("/employees", middleware.RequirePermission("employee:write"), handler.Create)
//	router.GET("/employees", middleware.RequirePermission("employee:read"), handler.List)
func RequirePermission(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// super_admin bypasses all checks
		roleName, _ := c.Get("role_name")
		if roleName == "super_admin" {
			c.Next()
			return
		}

		// Get permissions from JWT claims (set by JWTAuth middleware)
		permsVal, exists := c.Get("permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "no permissions found — contact your administrator",
			})
			return
		}

		permissions, ok := permsVal.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "invalid permissions format",
			})
			return
		}

		// Check if the required permission exists
		for _, p := range permissions {
			if p == required {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "you do not have permission to perform this action",
			"required": required,
		})
	}
}

// RequireAnyPermission checks if the user has at least one of the given permissions.
func RequireAnyPermission(perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName, _ := c.Get("role_name")
		if roleName == "super_admin" {
			c.Next()
			return
		}

		permsVal, exists := c.Get("permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "no permissions found — contact your administrator",
			})
			return
		}

		userPerms, ok := permsVal.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "invalid permissions format",
			})
			return
		}

		permSet := make(map[string]bool, len(userPerms))
		for _, p := range userPerms {
			permSet[p] = true
		}

		for _, required := range perms {
			if permSet[required] {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success":  false,
			"error":    "you do not have permission to perform this action",
			"required": perms,
		})
	}
}

// RequireRole checks if the authenticated user has a specific role name.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName, _ := c.Get("role_name")
		roleStr, ok := roleName.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "access denied",
			})
			return
		}

		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "insufficient role privileges",
		})
	}
}
