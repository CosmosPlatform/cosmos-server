package auth

import (
	"cosmos-server/api"
	"cosmos-server/pkg/auth"
	"github.com/gin-gonic/gin"
)

type handler struct {
	authService auth.Service
	translator  Translator
}

func AddAuthHandler(e *gin.RouterGroup, authService auth.Service) {
	handler := &handler{
		authService: authService,
		translator:  NewTranslator(),
	}

	authRoutes := e.Group("/auth")

	authRoutes.POST("/login", handler.handleLogin)
}

func (handler *handler) handleLogin(e *gin.Context) {
	var authenticateRequest api.AuthenticateRequest

	if err := e.ShouldBindJSON(&authenticateRequest); err != nil {
		e.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := authenticateRequest.Validate(); err != nil {
		e.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, token, err := handler.authService.Authenticate(e, authenticateRequest.Username, authenticateRequest.Password)
	if err != nil {
		e.JSON(401, gin.H{"error": "Authentication failed"})
		return
	}

	e.JSON(200, handler.translator.ToAuthenticateResponse(user, token))
}
