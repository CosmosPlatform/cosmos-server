package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/user"
	"github.com/gin-gonic/gin"
)

type handler struct {
	userService user.Service
	translator  Translator
}

func AddAdminUserHandler(e *gin.RouterGroup, userService user.Service) {
	handler := &handler{
		userService: userService,
		translator:  NewTranslator(),
	}

	e.POST("/users", handler.handleRegisterUser)
	e.POST("/adminUsers", handler.handleRegisterAdminUser)
}

func (handler *handler) handleRegisterUser(e *gin.Context) {
	var registerUserRequest api.RegisterUserRequest

	if err := e.ShouldBindJSON(&registerUserRequest); err != nil {
		e.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := registerUserRequest.Validate(); err != nil {
		e.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := handler.userService.RegisterRegularUser(e, registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)
	if err != nil {
		e.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	e.JSON(201, handler.translator.ToRegisterUserResponse(registerUserRequest.Username, registerUserRequest.Email, user.RegularUserRole))
}

func (handler *handler) handleRegisterAdminUser(e *gin.Context) {
	var registerUserRequest api.RegisterUserRequest

	if err := e.ShouldBindJSON(&registerUserRequest); err != nil {
		e.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := registerUserRequest.Validate(); err != nil {
		e.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := handler.userService.RegisterAdminUser(e, registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)
	if err != nil {
		e.JSON(500, gin.H{"error": "Failed to register admin user"})
		return
	}

	e.JSON(201, handler.translator.ToRegisterUserResponse(registerUserRequest.Username, registerUserRequest.Email, user.AdminUserRole))
}
