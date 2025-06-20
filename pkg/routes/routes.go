package routes

import (
	"cosmos-server/pkg/auth"
	authRoute "cosmos-server/pkg/routes/auth"
	healthcheckRoute "cosmos-server/pkg/routes/healthcheck"
	"github.com/gin-gonic/gin"
)

type HTTPRoutes struct {
	AuthService auth.Service
}

func NewHTTPRoutes(authService auth.Service) *HTTPRoutes {
	return &HTTPRoutes{
		AuthService: authService,
	}
}

func (r *HTTPRoutes) RegisterUnauthenticatedRoutes(e *gin.RouterGroup) {
	authRoute.AddAuthHandler(e, r.AuthService)
	healthcheckRoute.AddHealthcheckHandler(e)
}

func (r *HTTPRoutes) RegisterAdminAuthenticatedRoutes(e *gin.RouterGroup) {
	// TODO
}
