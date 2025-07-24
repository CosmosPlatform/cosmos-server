package routes

import (
	"cosmos-server/pkg/auth"
	"cosmos-server/pkg/log"
	authRoute "cosmos-server/pkg/routes/auth"
	healthcheckRoute "cosmos-server/pkg/routes/healthcheck"
	teamRoute "cosmos-server/pkg/routes/team"
	userRoute "cosmos-server/pkg/routes/user"
	"cosmos-server/pkg/team"
	"cosmos-server/pkg/user"
	"github.com/gin-gonic/gin"
)

type HTTPRoutes struct {
	AuthService auth.Service
	UserService user.Service
	TeamService team.Service
	Logger      log.Logger
}

func NewHTTPRoutes(authService auth.Service, userService user.Service, teamService team.Service, logger log.Logger) *HTTPRoutes {
	return &HTTPRoutes{
		AuthService: authService,
		UserService: userService,
		TeamService: teamService,
		Logger:      logger,
	}
}

func (r *HTTPRoutes) RegisterUnauthenticatedRoutes(e *gin.RouterGroup) {
	authRoute.AddAuthHandler(e, r.AuthService, r.Logger)
	healthcheckRoute.AddHealthcheckHandler(e)
}

func (r *HTTPRoutes) RegisterAdminAuthenticatedRoutes(e *gin.RouterGroup) {
	userRoute.AddAdminUserHandler(e, r.UserService, userRoute.NewTranslator(), r.Logger)
	teamRoute.AddTeamHandler(e, r.TeamService, teamRoute.NewTranslator())
}
