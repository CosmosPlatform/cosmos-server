package monitoring

import (
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/services/application"
	"cosmos-server/pkg/services/monitoring"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	monitoringService  monitoring.Service
	applicationService application.Service
	translator         Translator
	logger             log.Logger
}

func AddAuthenticatedMonitoringHandler(e *gin.RouterGroup, monitoringService monitoring.Service, applicationService application.Service, translator Translator, logger log.Logger) {
	handler := &handler{
		monitoringService:  monitoringService,
		applicationService: applicationService,
		translator:         translator,
		logger:             logger,
	}

	monitoringGroup := e.Group("/monitoring")

	monitoringGroup.POST("/update/:application", handler.handleUpdateApplicationMonitoring)
	monitoringGroup.GET("/interactions/:application", handler.handleGetApplicationInteractions)
}

func (handler *handler) handleUpdateApplicationMonitoring(e *gin.Context) {
	applicationName := e.Param("application")

	applicationToUpdate, err := handler.applicationService.GetApplication(e, applicationName)
	if err != nil {
		handler.logger.Errorf("Failed to retrieve application: %v", err)
		_ = e.Error(err)
		return
	}

	err = handler.monitoringService.UpdateApplicationInformation(e, applicationToUpdate)
	if err != nil {
		_ = e.Error(err)
		return
	}

	e.JSON(http.StatusNoContent, nil)
}

func (handler *handler) handleGetApplicationInteractions(e *gin.Context) {
	applicationName := e.Param("application")

	evaluatedApplication, err := handler.applicationService.GetApplication(e, applicationName)
	if err != nil {
		handler.logger.Errorf("Failed to retrieve evaluatedApplication: %v", err)
		_ = e.Error(err)
		return
	}

	interactions, err := handler.monitoringService.GetApplicationInteractions(e, evaluatedApplication.Name)
	if err != nil {
		handler.logger.Errorf("Failed to retrieve evaluatedApplication interactions: %v", err)
		_ = e.Error(err)
		return
	}

	interactions.MainApplication = evaluatedApplication

	e.JSON(200, handler.translator.ToGetApplicationInteractionsResponse(interactions))
}
