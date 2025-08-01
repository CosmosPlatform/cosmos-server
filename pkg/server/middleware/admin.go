package middleware

import (
	"cosmos-server/pkg/services/auth"
	"cosmos-server/pkg/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(auth.UserRoleContextKey)
		if !exists || role != user.AdminUserRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
		c.Next()
	}
}
