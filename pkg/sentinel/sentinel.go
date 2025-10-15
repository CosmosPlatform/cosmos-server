package sentinel

import (
	"context"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/services/application"
	"cosmos-server/pkg/services/monitoring"
)

type Sentinel struct {
	applicationService application.Service
	monitoringService  monitoring.Service
	newConfigChannel   <-chan model.SentinelSettings
	logger             log.Logger
}

func NewSentinel(logger log.Logger, applicationService application.Service, monitoringService monitoring.Service, newSettingsChannel <-chan model.SentinelSettings) *Sentinel {
	return &Sentinel{
		applicationService: applicationService,
		monitoringService:  monitoringService,
		newConfigChannel:   newSettingsChannel,
		logger:             logger,
	}
}

func (s *Sentinel) Start(ctx context.Context) {
	s.logger.Infof("Sentinel started")
}

