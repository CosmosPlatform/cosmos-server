package middleware

import (
	auth2 "cosmos-server/pkg/services/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleware(authService auth2.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(tokenHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with 'Bearer '"})
			return
		}
		token := strings.TrimPrefix(tokenHeader, bearerPrefix)

		validToken, err := authService.IsAuthenticated(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			return
		}

		claims, ok := validToken.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims from token"})
			return
		}

		if role, exists := claims[auth2.UserRoleClaimKey]; exists {
			c.Set(auth2.UserRoleContextKey, role)
		}

		c.Next()
	}
}
