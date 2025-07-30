package team

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/services/team"
	"fmt"
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
	e.POST("/teams", h.handleInsertTeam)
	e.DELETE("/teams", h.handleDeleteTeam)
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
	name := c.Query("name")
	if name == "" {
		_ = c.Error(errors.NewBadRequestError("team name is required"))
		return
	}

	err := h.teamService.DeleteTeam(c, name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(204)
}

func (h *handler) handleInsertTeam(c *gin.Context) {
	var teamRequest api.InsertTeamRequest

	if err := c.ShouldBindJSON(&teamRequest); err != nil {
		_ = c.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := teamRequest.Validate(); err != nil {
		_ = c.Error(err)
		return
	}

	teamModel := h.translator.ToModelTeam(teamRequest.Name, teamRequest.Description)
	err := h.teamService.InsertTeam(c, teamModel)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(201, h.translator.ToInsertTeamResponse(teamRequest.Name, teamRequest.Description))
}
