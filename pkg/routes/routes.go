package routes

import (
	"cosmos-server/pkg/auth"
	authRoute "cosmos-server/pkg/routes/auth"
	healthcheckRoute "cosmos-server/pkg/routes/healthcheck"
	userRoute "cosmos-server/pkg/routes/user"
	"cosmos-server/pkg/user"
	"github.com/gin-gonic/gin"
)

type HTTPRoutes struct {
	AuthService auth.Service
	UserService user.Service
}

func NewHTTPRoutes(authService auth.Service, userService user.Service) *HTTPRoutes {
	return &HTTPRoutes{
		AuthService: authService,
		UserService: userService,
	}
}

func (r *HTTPRoutes) RegisterUnauthenticatedRoutes(e *gin.RouterGroup) {
	authRoute.AddAuthHandler(e, r.AuthService)
	healthcheckRoute.AddHealthcheckHandler(e)
}

func (r *HTTPRoutes) RegisterAdminAuthenticatedRoutes(e *gin.RouterGroup) {
	userRoute.AddAdminUserHandler(e, r.UserService)
}
