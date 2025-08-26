package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/services/auth"
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
	usersGroup.DELETE("", handler.handleDeleteUser)
}

func AddAuthenticatedUserHandler(e *gin.RouterGroup, userService user.Service, translator Translator, logger log.Logger) {
	handler := &handler{
		userService: userService,
		translator:  translator,
		logger:      logger,
	}

	usersGroup := e.Group("/users")
	usersGroup.GET("/me", handler.handleGetCurrentUser)
}

func (handler *handler) handleRegisterUser(e *gin.Context) {
	var registerUserRequest api.RegisterUserRequest

	if err := e.ShouldBindJSON(&registerUserRequest); err != nil {
		handler.logger.Errorf("Failed to bind JSON for registration request: %v", err)
		_ = e.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := registerUserRequest.Validate(); err != nil {
		_ = e.Error(errors.NewBadRequestError(err.Error()))
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

func (handler *handler) handleDeleteUser(e *gin.Context) {
	email := e.Query("email")
	if email == "" {
		handler.logger.Errorf("Email query parameter is required for deleting a user")
		_ = e.Error(errors.NewBadRequestError("Email query parameter is required"))
		return
	}

	err := handler.userService.DeleteUser(e, email)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.Status(204)
}

func (handler *handler) handleGetCurrentUser(c *gin.Context) {
	userEmail, exists := c.Get(auth.UserEmailContextKey)
	if !exists {
		_ = c.Error(errors.NewUnauthorizedError("user email not found in context"))
		return
	}

	userModel, err := handler.userService.GetUserWithEmail(c, userEmail.(string))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, handler.translator.ToGetCurrentUserResponse(userModel))
}
