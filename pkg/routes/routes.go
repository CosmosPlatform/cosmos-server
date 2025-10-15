package routes

import (
	"cosmos-server/pkg/log"
	monitoringRoute "cosmos-server/pkg/routes/monitoring"
	"cosmos-server/pkg/services/application"
	"cosmos-server/pkg/services/auth"
	"cosmos-server/pkg/services/monitoring"
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
	MonitoringService  monitoring.Service
	Logger             log.Logger
}

func NewHTTPRoutes(authService auth.Service, userService user.Service, teamService team.Service, applicationService application.Service, monitoringService monitoring.Service, logger log.Logger) *HTTPRoutes {
	return &HTTPRoutes{
		AuthService:        authService,
		UserService:        userService,
		TeamService:        teamService,
		ApplicationService: applicationService,
		MonitoringService:  monitoringService,
		Logger:             logger,
	}
}

func (r *HTTPRoutes) RegisterUnauthenticatedRoutes(e *gin.RouterGroup) {
	authRoute.AddAuthHandler(e, r.AuthService, r.Logger)
	healthcheckRoute.AddHealthcheckHandler(e)
	monitoringRoute.AddAdminMonitoringHandler(e, r.MonitoringService, r.ApplicationService, monitoringRoute.NewTranslator(), r.Logger) // Mover luego

}

func (r *HTTPRoutes) RegisterAuthenticatedRoutes(e *gin.RouterGroup) {
	applicationRoute.AddAuthenticatedApplicationHandler(e, r.ApplicationService, r.MonitoringService, applicationRoute.NewTranslator(), r.Logger)
	userRoute.AddAuthenticatedUserHandler(e, r.UserService, userRoute.NewTranslator(), r.Logger)
	teamRoute.AddAuthenticatedTeamHandler(e, r.TeamService, teamRoute.NewTranslator())
	monitoringRoute.AddAuthenticatedMonitoringHandler(e, r.MonitoringService, r.ApplicationService, monitoringRoute.NewTranslator(), r.Logger)
}

func (r *HTTPRoutes) RegisterAdminAuthenticatedRoutes(e *gin.RouterGroup) {
	userRoute.AddAdminUserHandler(e, r.UserService, userRoute.NewTranslator(), r.Logger)
	teamRoute.AddAdminTeamHandler(e, r.TeamService, teamRoute.NewTranslator())
}
