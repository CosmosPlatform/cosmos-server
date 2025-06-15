package routes

import (
	"cosmos-server/pkg/auth"
	authRoute "cosmos-server/pkg/routes/auth"
	healthcheckRoute "cosmos-server/pkg/routes/healthcheck"
	"github.com/gin-gonic/gin"
)

type HTTPRoutes struct {
	authService auth.Service
}

func NewHTTPRoutes(authService auth.Service) *HTTPRoutes {
	return &HTTPRoutes{
		authService: authService,
	}
}

func (r *HTTPRoutes) RegisterUnauthenticatedRoutes(e *gin.RouterGroup) {
	authRoute.AddAuthHandler(e, r.authService)
	healthcheckRoute.AddHealthcheckHandler(e)
}
