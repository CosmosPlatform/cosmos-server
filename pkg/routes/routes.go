package routes

import (
	"cosmos-server/pkg/auth"
	"cosmos-server/pkg/log"
	authRoute "cosmos-server/pkg/routes/auth"
	healthcheckRoute "cosmos-server/pkg/routes/healthcheck"
	userRoute "cosmos-server/pkg/routes/user"
	"cosmos-server/pkg/user"
	"github.com/gin-gonic/gin"
)

type HTTPRoutes struct {
	AuthService auth.Service
	UserService user.Service
	Logger      log.Logger
}

func NewHTTPRoutes(authService auth.Service, userService user.Service, logger log.Logger) *HTTPRoutes {
	return &HTTPRoutes{
		AuthService: authService,
		UserService: userService,
		Logger:      logger,
	}
}

func (r *HTTPRoutes) RegisterUnauthenticatedRoutes(e *gin.RouterGroup) {
	authRoute.AddAuthHandler(e, r.AuthService, r.Logger)
	healthcheckRoute.AddHealthcheckHandler(e)
}

func (r *HTTPRoutes) RegisterAdminAuthenticatedRoutes(e *gin.RouterGroup) {
	userRoute.AddAdminUserHandler(e, r.UserService)
}
