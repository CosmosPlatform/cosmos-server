package auth

import (
	"cosmos-server/pkg/auth"
	"github.com/gin-gonic/gin"
)

type handler struct {
	authService auth.Service
}

func AddAuthHandler(e *gin.RouterGroup, authService auth.Service) {
	handler := &handler{
		authService: authService,
	}

	authRoutes := e.Group("/auth")

	authRoutes.POST("/login", handler.handleLogin)
}

func (handler *handler) handleLogin(e *gin.Context) {

}
