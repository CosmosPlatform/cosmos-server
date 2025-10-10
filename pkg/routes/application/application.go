package application

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/services/application"
	"cosmos-server/pkg/services/monitoring"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	applicationService application.Service
	monitoringService  monitoring.Service
	translator         Translator
	logger             log.Logger
}

func AddAuthenticatedApplicationHandler(e *gin.RouterGroup, applicationService application.Service, monitoringService monitoring.Service, translator Translator, logger log.Logger) {
	handler := &handler{
		applicationService: applicationService,
		monitoringService:  monitoringService,
		translator:         translator,
		logger:             logger,
	}

	applicationsGroup := e.Group("/applications")

	applicationsGroup.GET("/:application", handler.handleGetApplication)
	applicationsGroup.GET("", handler.handleGetApplications)
	applicationsGroup.GET("/team/:teamName", handler.handleGetApplicationByTeam)
	applicationsGroup.POST("", handler.handleCreateApplication)
	applicationsGroup.PUT("/:application", handler.handleUpdateApplication)
	applicationsGroup.DELETE("/:application", handler.handleDeleteApplication)
}

func (handler *handler) handleCreateApplication(e *gin.Context) {
	var createApplicationRequest api.CreateApplicationRequest

	if err := e.ShouldBindJSON(&createApplicationRequest); err != nil {
		handler.logger.Errorf("Failed to bind JSON for registration request: %v", err)
		_ = e.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := createApplicationRequest.Validate(); err != nil {
		_ = e.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	gitInformation := handler.translator.ToGitInformationModel(createApplicationRequest.GitInformation)
	monitoringInformation := handler.translator.ToMonitoringInformationModel(createApplicationRequest.MonitoringInformation)

	err := handler.applicationService.AddApplication(e, createApplicationRequest.Name, createApplicationRequest.Description, createApplicationRequest.Team, gitInformation, monitoringInformation)
	if err != nil {
		_ = e.Error(err)
		return
	}

	if gitInformation != nil {
		app, err := handler.applicationService.GetApplication(e, createApplicationRequest.Name)
		if err != nil {
			handler.logger.Errorf("Failed to retrieve application after creation: %v", err)
			_ = e.Error(err)
			return
		}

		err = handler.monitoringService.UpdateApplicationInformation(e, app)
		if err != nil {
			handler.logger.Errorf("Failed to update application information after creation: %v", err)
			_ = e.Error(err)
			return
		}
	}

	e.JSON(http.StatusCreated, handler.translator.ToCreateApplicationResponse(createApplicationRequest.Name, createApplicationRequest.Description, createApplicationRequest.Team, gitInformation))
}

func (handler *handler) handleGetApplication(e *gin.Context) {
	applicationName := e.Param("application")
	if applicationName == "" {
		_ = e.Error(errors.NewBadRequestError("application name missing"))
		return
	}

	app, err := handler.applicationService.GetApplication(e, applicationName)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusOK, handler.translator.ToGetApplicationResponse(app))
}

func (handler *handler) handleGetApplications(e *gin.Context) {
	name := e.Query("name")

	applications, err := handler.applicationService.GetApplicationsWithFilter(e, name)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusOK, handler.translator.ToGetApplicationsResponse(applications))
}

func (handler *handler) handleGetApplicationByTeam(e *gin.Context) {
	teamName := e.Param("teamName")
	if teamName == "" {
		_ = e.Error(errors.NewBadRequestError("team name missing"))
		return
	}

	applications, err := handler.applicationService.GetApplicationsByTeam(e, teamName)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusOK, handler.translator.ToGetApplicationsResponse(applications))
}

func (handler *handler) handleDeleteApplication(e *gin.Context) {
	applicationName := e.Param("application")
	if applicationName == "" {
		_ = e.Error(errors.NewBadRequestError("application name missing"))
		return
	}

	err := handler.applicationService.DeleteApplication(e, applicationName)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.Status(http.StatusNoContent)
}

func (handler *handler) handleUpdateApplication(e *gin.Context) {
	applicationName := e.Param("application")
	if applicationName == "" {
		_ = e.Error(errors.NewBadRequestError("application name missing"))
		return
	}

	var updateRequest api.UpdateApplicationRequest
	if err := e.ShouldBindJSON(&updateRequest); err != nil {
		handler.logger.Errorf("Failed to bind JSON for update request: %v", err)
		_ = e.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	if err := updateRequest.Validate(); err != nil {
		_ = e.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	updateData := &model.ApplicationUpdate{
		Name:        updateRequest.Name,
		Description: updateRequest.Description,
		Team:        updateRequest.Team,
	}

	if updateRequest.GitInformation != nil {
		updateData.GitInformation = &model.GitInformation{
			Provider:         updateRequest.GitInformation.Provider,
			RepositoryOwner:  updateRequest.GitInformation.RepositoryOwner,
			RepositoryName:   updateRequest.GitInformation.RepositoryName,
			RepositoryBranch: updateRequest.GitInformation.RepositoryBranch,
		}
		if updateRequest.MonitoringInformation != nil {
			updateData.MonitoringInformation = &model.MonitoringInformation{
				HasOpenApi:     updateRequest.MonitoringInformation.HasOpenAPI,
				HasOpenClient:  updateRequest.MonitoringInformation.HasOpenClient,
				OpenApiPath:    updateRequest.MonitoringInformation.OpenAPIPath,
				OpenClientPath: updateRequest.MonitoringInformation.OpenClientPath,
			}
		}
	}

	updatedApp, err := handler.applicationService.UpdateApplication(e, applicationName, updateData)
	if err != nil {
		_ = e.Error(err)
		return
	}

	if updateData.GitInformation != nil {
		err = handler.monitoringService.UpdateApplicationInformation(e, updatedApp)
		if err != nil {
			handler.logger.Errorf("Failed to update application information after update: %v", err)
			_ = e.Error(err)
			return
		}
	}

	e.JSON(http.StatusOK, handler.translator.ToUpdateApplicationResponse(updatedApp))
}
