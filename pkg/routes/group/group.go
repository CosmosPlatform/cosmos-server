package group

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/services/group"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	groupService group.Service
	translator   Translator
	logger       log.Logger
}

func AddAuthenticatedGroupHandler(e *gin.RouterGroup, groupService group.Service, translator Translator, logger log.Logger) {
	handler := &handler{
		groupService: groupService,
		translator:   translator,
		logger:       logger,
	}

	groupsGroup := e.Group("/groups")

	groupsGroup.POST("", handler.handleCreateGroup)
	groupsGroup.GET("", handler.handleListGroups)
	groupsGroup.GET("/:groupName", handler.handleGetGroup)
}

func (handler *handler) handleCreateGroup(e *gin.Context) {
	var createGroupRequest api.CreateGroupRequest
	if err := e.ShouldBindJSON(&createGroupRequest); err != nil {
		handler.logger.Errorf("Failed to bind JSON for create group request: %v", err)
		_ = e.Error(err)
		return
	}

	err := createGroupRequest.Validate()
	if err != nil {
		_ = e.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	err = handler.groupService.CreateGroup(e, createGroupRequest.Name, createGroupRequest.Description, createGroupRequest.Members)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.Status(http.StatusCreated)
}

func (handler *handler) handleListGroups(e *gin.Context) {
	groups, err := handler.groupService.GetAllGroups(e)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusOK, handler.translator.ToGetGroupsResponse(groups))
}

func (handler *handler) handleGetGroup(e *gin.Context) {
	groupName := e.Param("groupName")

	groupObj, err := handler.groupService.GetGroupByName(e, groupName)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusOK, handler.translator.ToGetGroupResponse(groupObj))
}
