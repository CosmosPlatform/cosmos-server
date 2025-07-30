package middleware

import (
	"cosmos-server/pkg/services/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleware(authService auth.Service) gin.HandlerFunc {
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

		if role, exists := claims[auth.UserRoleClaimKey]; exists {
			c.Set(auth.UserRoleContextKey, role)
		}

		if email, exists := claims[auth.UserEmailClaimKey]; exists {
			c.Set(auth.UserEmailContextKey, email)
		}

		c.Next()
	}
}
