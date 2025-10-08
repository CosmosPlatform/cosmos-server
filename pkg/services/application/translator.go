package application

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

type Translator interface {
	ToApplicationModel(applicationObj *obj.Application) *model.Application
	ToApplicationModels(applicationObjs []*obj.Application) []*model.Application
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToApplicationModel(applicationObj *obj.Application) *model.Application {
	modelApplication := &model.Application{
		Name:                  applicationObj.Name,
		Description:           applicationObj.Description,
		Team:                  t.ToModelTeam(applicationObj.Team),
		MonitoringInformation: t.ToModelMonitoringInformation(applicationObj),
	}

	if applicationObj.GitProvider != "" || applicationObj.GitRepositoryName != "" || applicationObj.GitRepositoryOwner != "" || applicationObj.GitRepositoryBranch != "" {
		modelApplication.GitInformation = &model.GitInformation{
			Provider:         applicationObj.GitProvider,
			RepositoryOwner:  applicationObj.GitRepositoryOwner,
			RepositoryName:   applicationObj.GitRepositoryName,
			RepositoryBranch: applicationObj.GitRepositoryBranch,
		}
	}

	return modelApplication
}

func (t *translator) ToModelTeam(teamObj *obj.Team) *model.Team {
	if teamObj == nil {
		return nil
	}
	return &model.Team{
		Name:        teamObj.Name,
		Description: teamObj.Description,
	}
}

func (t *translator) ToApplicationModels(applicationObjs []*obj.Application) []*model.Application {
	applicationModels := make([]*model.Application, 0, len(applicationObjs))
	for _, applicationObj := range applicationObjs {
		applicationModels = append(applicationModels, t.ToApplicationModel(applicationObj))
	}
	return applicationModels
}

func (t *translator) ToModelMonitoringInformation(applicationObj *obj.Application) *model.MonitoringInformation {
	if applicationObj == nil {
		return nil
	}

	return &model.MonitoringInformation{
		DependenciesSha: applicationObj.DependenciesSha,
		OpenAPISha:      applicationObj.OpenAPISha,
	}
}
