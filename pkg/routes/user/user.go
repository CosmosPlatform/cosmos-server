package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/services/user"
	"fmt"
	"github.com/gin-gonic/gin"
)

type handler struct {
	userService user.Service
	translator  Translator
	logger      log.Logger
}

func AddAdminUserHandler(e *gin.RouterGroup, userService user.Service, translator Translator, logger log.Logger) {
	handler := &handler{
		userService: userService,
		translator:  translator,
		logger:      logger,
	}

	usersGroup := e.Group("/users")

	usersGroup.GET("", handler.handleGetUsers)
	usersGroup.POST("", handler.handleRegisterUser)
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

	err := handler.userService.RegisterUser(e, registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password, registerUserRequest.Role)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(201, handler.translator.ToRegisterUserResponse(registerUserRequest.Username, registerUserRequest.Email, user.RegularUserRole))
}

func (handler *handler) handleGetUsers(e *gin.Context) {
	users, err := handler.userService.GetUsers(e)
	if err != nil {
		handler.logger.Errorf("Failed to get users: %v", err)
		_ = e.Error(err)
		return
	}

	e.JSON(200, handler.translator.ToGetUsersResponse(users))
}
