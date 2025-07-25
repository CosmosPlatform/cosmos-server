package team

import (
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/services/team"
	"cosmos-server/pkg/storage"
	"fmt"
	"github.com/gin-gonic/gin"

	errorUtils "errors"
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
	e.DELETE("/teams/:name", h.handleDeleteTeam)
}

func (h *handler) handleGetTeams(c *gin.Context) {
	teams, err := h.teamService.GetAllTeams(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(200, h.translator.ToGetTeamsResponse(teams))
}

func (h *handler) handleDeleteTeam(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		_ = c.Error(errors.NewBadRequestError("team name is required"))
		return
	}

	err := h.teamService.DeleteTeam(c, name)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			_ = c.Error(errors.NewNotFoundError("team not found"))
			return
		}
		_ = c.Error(errors.NewInternalServerError(fmt.Sprintf("failed to delete team: %v", err)))
		return
	}

	c.Status(204)
}
