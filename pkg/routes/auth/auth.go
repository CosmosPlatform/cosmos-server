package auth

import (
	"cosmos-server/api"
	"cosmos-server/pkg/auth"
	"cosmos-server/pkg/errors"
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
	var authenticateRequest api.AuthenticateRequest

	if err := e.ShouldBindJSON(&authenticateRequest); err != nil {
		handler.logger.Errorf("Failed to bind JSON for authentication request: %v", err)
		_ = e.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := authenticateRequest.Validate(); err != nil {
		handler.logger.Errorf("Validation error for authentication request: %v", err)
		_ = e.Error(err)
		return
	}

	user, token, err := handler.authService.Authenticate(e, authenticateRequest.Email, authenticateRequest.Password)
	if err != nil {
		handler.logger.Errorf("Authentication failed for user %s: %v", authenticateRequest.Email, err)
		_ = e.Error(err)
		return
	}

	e.JSON(200, handler.translator.ToAuthenticateResponse(user, token))
}
