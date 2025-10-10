package application

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToCreateApplicationResponse(name, description, team string, gitInformation *model.GitInformation) *api.CreateApplicationResponse
	ToGetApplicationResponse(applicationObj *model.Application) *api.GetApplicationResponse
	ToGetApplicationsResponse(applicationObj []*model.Application) *api.GetApplicationsResponse
	ToUpdateApplicationResponse(app *model.Application) *api.UpdateApplicationResponse
	ToGitInformationModel(gitInfo *api.GitInformation) *model.GitInformation
	ToMonitoringInformationModel(monitoringInfo *api.MonitoringInformation) *model.MonitoringInformation
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToCreateApplicationResponse(name, description, team string, gitInformation *model.GitInformation) *api.CreateApplicationResponse {
	return &api.CreateApplicationResponse{
		Application: &api.Application{
			Name:           name,
			Description:    description,
			Team:           &api.Team{Name: team},
			GitInformation: t.ToApiGitInformation(gitInformation),
		},
	}
}

func (t *translator) ToGetApplicationResponse(applicationModel *model.Application) *api.GetApplicationResponse {
	return &api.GetApplicationResponse{
		Application: t.ToApplicationApi(applicationModel),
	}
}

func (t *translator) ToApiTeam(teamModel *model.Team) *api.Team {
	if teamModel == nil {
		return nil
	}
	return &api.Team{
		Name:        teamModel.Name,
		Description: teamModel.Description,
	}
}

func (t *translator) ToApiGitInformation(gitInfo *model.GitInformation) *api.GitInformation {
	if gitInfo == nil {
		return nil
	}
	return &api.GitInformation{
		Provider:         gitInfo.Provider,
		RepositoryOwner:  gitInfo.RepositoryOwner,
		RepositoryName:   gitInfo.RepositoryName,
		RepositoryBranch: gitInfo.RepositoryBranch,
	}
}

func (t *translator) ToGetApplicationsResponse(applicationModels []*model.Application) *api.GetApplicationsResponse {
	applications := make([]*api.Application, 0)
	for _, applicationModel := range applicationModels {
		applications = append(applications, t.ToApplicationApi(applicationModel))
	}
	return &api.GetApplicationsResponse{
		Applications: applications,
	}
}

func (t *translator) ToApplicationApi(applicationModel *model.Application) *api.Application {
	return &api.Application{
		Name:                  applicationModel.Name,
		Description:           applicationModel.Description,
		Team:                  t.ToApiTeam(applicationModel.Team),
		GitInformation:        t.ToApiGitInformation(applicationModel.GitInformation),
		MonitoringInformation: t.ToMonitoringInformationApi(applicationModel.MonitoringInformation),
	}
}

func (t *translator) ToUpdateApplicationResponse(app *model.Application) *api.UpdateApplicationResponse {
	return &api.UpdateApplicationResponse{
		Application: t.ToApplicationApi(app),
	}
}

func (t *translator) ToGitInformationModel(gitInfo *api.GitInformation) *model.GitInformation {
	if gitInfo == nil {
		return nil
	}
	return &model.GitInformation{
		Provider:         gitInfo.Provider,
		RepositoryOwner:  gitInfo.RepositoryOwner,
		RepositoryName:   gitInfo.RepositoryName,
		RepositoryBranch: gitInfo.RepositoryBranch,
	}
}

func (t *translator) ToMonitoringInformationModel(monitoringInfo *api.MonitoringInformation) *model.MonitoringInformation {
	if monitoringInfo == nil {
		return nil
	}

	return &model.MonitoringInformation{
		HasOpenApi:     monitoringInfo.HasOpenApi,
		OpenApiPath:    monitoringInfo.OpenApiPath,
		HasOpenClient:  monitoringInfo.HasOpenClient,
		OpenClientPath: monitoringInfo.OpenClientPath,
	}
}

func (t *translator) ToMonitoringInformationApi(monitoringInfo *model.MonitoringInformation) *api.MonitoringInformation {
	if monitoringInfo == nil {
		return nil
	}

	return &api.MonitoringInformation{
		HasOpenApi:     monitoringInfo.HasOpenApi,
		OpenApiPath:    monitoringInfo.OpenApiPath,
		HasOpenClient:  monitoringInfo.HasOpenClient,
		OpenClientPath: monitoringInfo.OpenClientPath,
	}
}
