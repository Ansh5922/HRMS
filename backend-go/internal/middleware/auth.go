package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-org/hrms-backend/internal/config"
	pkgjwt "github.com/your-org/hrms-backend/pkg/jwt"
)

// JWTAuth validates the JWT token and injects user_id, org_id, role_id,
// role_name, and permissions into the Gin context.
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := pkgjwt.ValidateToken(tokenStr, cfg.AppSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid or expired token"})
			return
		}

		// Inject into context for downstream handlers & middleware
		c.Set("user_id", claims.UserID)
		c.Set("org_id", claims.OrgID)
		c.Set("role_id", claims.RoleID)
		c.Set("role_name", claims.RoleName)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}
