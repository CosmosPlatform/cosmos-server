package routes

import (
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/services/application"
	"cosmos-server/pkg/services/auth"
	"cosmos-server/pkg/services/team"
	"cosmos-server/pkg/services/user"
	"github.com/gin-gonic/gin"

	applicationRoute "cosmos-server/pkg/routes/application"
	authRoute "cosmos-server/pkg/routes/auth"
	healthcheckRoute "cosmos-server/pkg/routes/healthcheck"
	teamRoute "cosmos-server/pkg/routes/team"
	userRoute "cosmos-server/pkg/routes/user"
)

type HTTPRoutes struct {
	AuthService        auth.Service
	UserService        user.Service
	TeamService        team.Service
	ApplicationService application.Service
	Logger             log.Logger
}

func NewHTTPRoutes(authService auth.Service, userService user.Service, teamService team.Service, applicationService application.Service, logger log.Logger) *HTTPRoutes {
	return &HTTPRoutes{
		AuthService:        authService,
		UserService:        userService,
		TeamService:        teamService,
		ApplicationService: applicationService,
		Logger:             logger,
	}
}

func (r *HTTPRoutes) RegisterUnauthenticatedRoutes(e *gin.RouterGroup) {
	authRoute.AddAuthHandler(e, r.AuthService, r.Logger)
	healthcheckRoute.AddHealthcheckHandler(e)
}

func (r *HTTPRoutes) RegisterAuthenticatedRoutes(e *gin.RouterGroup) {
	applicationRoute.AddApplicationHandler(e, r.ApplicationService, applicationRoute.NewTranslator(), r.Logger)
}

func (r *HTTPRoutes) RegisterAdminAuthenticatedRoutes(e *gin.RouterGroup) {
	userRoute.AddAdminUserHandler(e, r.UserService, userRoute.NewTranslator(), r.Logger)
	teamRoute.AddAdminTeamHandler(e, r.TeamService, teamRoute.NewTranslator())
}
