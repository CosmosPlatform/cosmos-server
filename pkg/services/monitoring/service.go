package monitoring

import (
	"context"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"encoding/json"
	"strings"
)

type Service interface {
	UpdateApplicationInformation(ctx context.Context, application *model.Application) error
}

type monitoringService struct {
	storageService storage.Service
	gitService     GitService
	translator     Translator
	logger         log.Logger
}

func NewMonitoringService(storageService storage.Service, gitService GitService, translator Translator, logger log.Logger) Service {
	return &monitoringService{
		storageService: storageService,
		gitService:     gitService,
		translator:     translator,
		logger:         logger,
	}
}

func (s *monitoringService) UpdateApplicationInformation(ctx context.Context, application *model.Application) error {
	if application.GitInformation == nil {
		s.logger.Infof("No git information for application %s, skipping monitoring update", application.Name)
		return nil // Could be an error because there is nothing to update.
	}

	rawOpenClientDefinition, err := s.gitService.GetFileWithContent(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, "docs/openclient.json")
	if err != nil {
		s.logger.Errorf("Failed to get openclient.json for application %s: %v", application.Name, err)
		return err
	}

	var openClientDef model.OpenClientSpecification
	decoder := json.NewDecoder(strings.NewReader(rawOpenClientDefinition.Content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&openClientDef); err != nil {
		s.logger.Errorf("Failed to unmarshal openclient.json for application %s: %v", application.Name, err)
		return err
	}

	if err := openClientDef.Validate(); err != nil {
		s.logger.Errorf("Invalid openclient.json for application %s: %v", application.Name, err)
		return err
	}

	for _, dependency := range openClientDef.Dependencies {
		s.logger.Infof("Application %s has dependency with reasons: %v", application.Name, dependency.Reasons)
		for path, methods := range dependency.Endpoints {
			for method, details := range methods {
				s.logger.Infof("Application %s has endpoint %s %s with reasons: %v", application.Name, method, path, details.Reasons)
			}
		}
	}

	return nil
}
