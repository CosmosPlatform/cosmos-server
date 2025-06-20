package server

import (
	"cosmos-server/pkg/auth"
	"cosmos-server/pkg/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func authMiddleware(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		token, err := authService.IsAuthenticated(tokenHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims from token"})
			return
		}

		if userID, exists := claims[auth.UserIDClaimKey]; exists {
			c.Set(auth.UserIDContextKey, userID)
		}

		if role, exists := claims[auth.UserRoleClaimKey]; exists {
			c.Set(auth.UserRoleContextKey, role)
		}

		c.Next()
	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(auth.UserRoleContextKey)
		if !exists || role != user.AdminUserRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
		c.Next()
	}
}
