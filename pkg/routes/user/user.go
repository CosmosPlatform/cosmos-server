package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/user"
	"fmt"
	"github.com/gin-gonic/gin"
)

type handler struct {
	userService user.Service
	translator  Translator
	logger      log.Logger
}

func AddAdminUserHandler(e *gin.RouterGroup, userService user.Service, logger log.Logger) {
	handler := &handler{
		userService: userService,
		translator:  NewTranslator(),
		logger:      logger,
	}

	e.POST("/users", handler.handleRegisterUser)
	e.POST("/adminUsers", handler.handleRegisterAdminUser)
}

func (handler *handler) handleRegisterUser(e *gin.Context) {
	var registerUserRequest api.RegisterUserRequest

	if err := e.ShouldBindJSON(&registerUserRequest); err != nil {
		handler.logger.Errorf("Failed to bind JSON for registration request: %v", err)
		_ = e.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := registerUserRequest.Validate(); err != nil {
		_ = e.Error(err)
		return
	}

	err := handler.userService.RegisterRegularUser(e, registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(201, handler.translator.ToRegisterUserResponse(registerUserRequest.Username, registerUserRequest.Email, user.RegularUserRole))
}

func (handler *handler) handleRegisterAdminUser(e *gin.Context) {
	var registerUserRequest api.RegisterUserRequest

	if err := e.ShouldBindJSON(&registerUserRequest); err != nil {

		_ = e.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := registerUserRequest.Validate(); err != nil {
		e.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := handler.userService.RegisterAdminUser(e, registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(201, handler.translator.ToRegisterUserResponse(registerUserRequest.Username, registerUserRequest.Email, user.AdminUserRole))
}
