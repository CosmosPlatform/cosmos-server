package team

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/services/team"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	teamService team.Service
	translator  Translator
}

func AddAdminTeamHandler(e *gin.RouterGroup, teamService team.Service, translator Translator) {
	h := &handler{
		teamService: teamService,
		translator:  translator,
	}

	e.GET("/teams", h.handleGetTeams)
	e.POST("/teams", h.handleCreateTeam)
	e.DELETE("/teams/:teamName", h.handleDeleteTeam)

	e.POST("/teams/:teamName/members", h.handleAddUserToTeam)
	e.DELETE("/teams/:teamName/members", h.handleRemoveUserFromTeam)
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
	name := c.Param("teamName")
	if name == "" {
		_ = c.Error(errors.NewBadRequestError("team name is required"))
		return
	}

	err := h.teamService.DeleteTeam(c, name)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handler) handleCreateTeam(c *gin.Context) {
	var teamRequest api.CreateTeamRequest

	if err := c.ShouldBindJSON(&teamRequest); err != nil {
		_ = c.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := teamRequest.Validate(); err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	teamModel := h.translator.ToModelTeam(teamRequest.Name, teamRequest.Description)
	err := h.teamService.InsertTeam(c, teamModel)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, h.translator.ToInsertTeamResponse(teamRequest.Name, teamRequest.Description))
}

func (h *handler) handleAddUserToTeam(c *gin.Context) {
	teamName := c.Param("teamName")
	if teamName == "" {
		_ = c.Error(errors.NewBadRequestError("team name is required"))
		return
	}

	var addUserRequest api.AddUserToTeamRequest
	if err := c.ShouldBindJSON(&addUserRequest); err != nil {
		_ = c.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := addUserRequest.Validate(); err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	err := h.teamService.AddUserToTeam(c, addUserRequest.Email, teamName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handler) handleRemoveUserFromTeam(c *gin.Context) {
	teamName := c.Param("teamName")
	if teamName == "" {
		_ = c.Error(errors.NewBadRequestError("team name is required"))
		return
	}

	email := c.Query("email")
	if email == "" {
		_ = c.Error(errors.NewBadRequestError("email query parameter is required"))
		return
	}

	_, err := h.teamService.GetTeamByName(c, teamName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.teamService.RemoveUserFromTeam(c, email)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
