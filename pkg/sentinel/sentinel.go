package sentinel

import (
	"context"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/services/application"
	"cosmos-server/pkg/services/monitoring"
	"sync"
	"time"
)

type Sentinel struct {
	applicationService application.Service
	monitoringService  monitoring.Service
	newConfigChannel   <-chan model.SentinelSettings
	workerCount        int
	jobsChan           chan *model.Application
	logger             log.Logger
}

func NewSentinel(logger log.Logger, applicationService application.Service, monitoringService monitoring.Service, newSettingsChannel <-chan model.SentinelSettings, workerCount int) *Sentinel {
	return &Sentinel{
		applicationService: applicationService,
		monitoringService:  monitoringService,
		newConfigChannel:   newSettingsChannel,
		workerCount:        workerCount,
		jobsChan:           make(chan *model.Application, 200),
		logger:             logger,
	}
}

func (s *Sentinel) Start(ctx context.Context, fallbackConfig *model.SentinelSettings) {
	s.logger.Infof("Sentinel starting with %d workers...", s.workerCount)

	var wg sync.WaitGroup
	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)
		go s.worker(ctx, i, &wg)
	}

	initialSettings, err := s.monitoringService.GetSentinelSettings(ctx)
	if err != nil {
		s.logger.Errorf("Failed to get initial sentinel settings: %v", err)
		initialSettings = fallbackConfig
	}

	ticker := time.NewTicker(time.Duration(initialSettings.Interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			s.logger.Infof("Sentinel activated: checking applications for updates")
			err := s.monitorApplications(ctx)
			if err != nil {
				s.logger.Errorf("Error checking applications: %v", err)
			}
		case newSettings := <-s.newConfigChannel:
			s.logger.Infof("Received new sentinel settings: %+v", newSettings)
			ticker.Stop()
			if newSettings.Enabled {
				s.logger.Infof("Sentinel reconfigured: interval set to %d seconds", newSettings.Interval)
				ticker = time.NewTicker(time.Duration(newSettings.Interval) * time.Second)
			} else {
				s.logger.Infof("Sentinel disabled")
				ticker = &time.Ticker{}
				ticker.C = nil
			}
		case <-ctx.Done():
			s.logger.Infof("Sentinel stopping...")
			close(s.jobsChan)
			wg.Wait()
			s.logger.Infof("Sentinel stopped")
			return
		}
	}
}

func (s *Sentinel) monitorApplications(ctx context.Context) error {
	applications, err := s.applicationService.GetApplicationsToMonitor(ctx)
	if err != nil {
		return err
	}

	for _, app := range applications {
		select {
		case s.jobsChan <- app:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (s *Sentinel) worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case app, ok := <-s.jobsChan:
			if !ok {
				s.logger.Infow("Jobs channel closed, exiting", "worker", id)
				return
			}
			s.logger.Infow("Monitoring application", "worker", id, "application", app.Name)
			s.monitorApplication(ctx, app, id)
		case <-ctx.Done():
			s.logger.Infow("Context cancelled", "worker", id)
			return
		}
	}
}

func (s *Sentinel) monitorApplication(ctx context.Context, app *model.Application, workerID int) {
	if err := s.monitoringService.UpdateApplicationDependencies(ctx, app); err != nil {
		s.logger.Errorf("Worker %d: Failed to update dependencies for application %s: %v", workerID, app.Name, err)
	}

	if err := s.monitoringService.UpdateApplicationOpenAPISpecification(ctx, app); err != nil {
		s.logger.Errorf("Worker %d: Failed to update OpenAPI specification for application %s: %v", workerID, app.Name, err)
	}
}
