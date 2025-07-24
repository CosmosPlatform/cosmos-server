package team

import (
	"cosmos-server/pkg/team"
	"github.com/gin-gonic/gin"
)

type handler struct {
	teamService team.Service
	translator  Translator
}

func AddTeamHandler(e *gin.RouterGroup, teamService team.Service, translator Translator) {
	h := &handler{
		teamService: teamService,
		translator:  translator,
	}

	e.GET("/teams", h.handleGetTeams)
}

func (h *handler) handleGetTeams(c *gin.Context) {
	teams, err := h.teamService.GetAllTeams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, h.translator.ToGetTeamsResponse(teams))
}
