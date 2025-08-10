package application

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/services/application"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	applicationService application.Service
	translator         Translator
	logger             log.Logger
}

func AddApplicationHandler(e *gin.RouterGroup, applicationService application.Service, translator Translator, logger log.Logger) {
	handler := &handler{
		applicationService: applicationService,
		translator:         translator,
		logger:             logger,
	}

	applicationsGroup := e.Group("/applications")

	applicationsGroup.GET("/:application", handler.handleGetApplication)
	applicationsGroup.POST("", handler.handleCreateApplication)
	//applicationsGroup.DELETE("/:id", handler.handleDeleteApplication)
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

	err := handler.applicationService.AddApplication(e, createApplicationRequest.Name, createApplicationRequest.Description, createApplicationRequest.Team)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusCreated, handler.translator.ToCreateApplicationResponse(createApplicationRequest.Name, createApplicationRequest.Description, createApplicationRequest.Team))
}

func (handler *handler) handleGetApplication(e *gin.Context) {

	applicationName := e.Param("application")
	if applicationName == "" {
		_ = e.Error(errors.NewBadRequestError("application name missing"))
	}

	app, err := handler.applicationService.GetApplication(e, applicationName)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusOK, handler.translator.ToGetApplicationResponse(app))

}
