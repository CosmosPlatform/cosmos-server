package auth

import (
	"cosmos-server/api"
	"cosmos-server/pkg/auth"
	"cosmos-server/pkg/log"
	"fmt"
	"github.com/gin-gonic/gin"
)

type handler struct {
	authService auth.Service
	translator  Translator
	logger      log.Logger
}

func AddAuthHandler(e *gin.RouterGroup, authService auth.Service, logger log.Logger) {
	handler := &handler{
		authService: authService,
		translator:  NewTranslator(),
		logger:      logger,
	}

	authRoutes := e.Group("/auth")

	authRoutes.POST("/login", handler.handleLogin)
}

func (handler *handler) handleLogin(e *gin.Context) {
	fmt.Printf("Handling login request\n")
	var authenticateRequest api.AuthenticateRequest

	if err := e.ShouldBindJSON(&authenticateRequest); err != nil {
		handler.logger.Errorf("Failed to bind JSON for authentication request: %v", err)
		e.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := authenticateRequest.Validate(); err != nil {
		handler.logger.Errorf("Validation error for authentication request: %v", err)
		e.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, token, err := handler.authService.Authenticate(e, authenticateRequest.Email, authenticateRequest.Password)
	if err != nil {
		handler.logger.Errorf("Authentication failed for user %s: %v", authenticateRequest.Email, err)
		e.JSON(401, gin.H{"error": "Authentication failed"})
		return
	}

	e.JSON(200, handler.translator.ToAuthenticateResponse(user, token))
}
