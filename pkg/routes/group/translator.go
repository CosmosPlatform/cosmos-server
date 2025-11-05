package group

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToGetGroupsResponse(modelGroups []*model.Group) api.GetGroupsResponse
	ToGetGroupResponse(modelGroup *model.Group) api.GetGroupResponse
}

type translator struct {
}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToGetGroupsResponse(modelGroups []*model.Group) api.GetGroupsResponse {
	apiReducedGroups := make([]*api.GroupReduced, 0, len(modelGroups))

	for _, modelGroup := range modelGroups {
		apiReducedGroup := api.GroupReduced{
			Name:        modelGroup.Name,
			Description: modelGroup.Description,
		}
		apiReducedGroups = append(apiReducedGroups, &apiReducedGroup)
	}

	return api.GetGroupsResponse{
		Groups: apiReducedGroups,
	}
}

func (t *translator) ToGetGroupResponse(modelGroup *model.Group) api.GetGroupResponse {
	if modelGroup == nil {
		return api.GetGroupResponse{}
	}

	return api.GetGroupResponse{
		Group: &api.Group{
			Name:        modelGroup.Name,
			Description: modelGroup.Description,
			Members:     t.ToGetApplicationsApi(modelGroup.Members),
		},
	}
}

func (t *translator) ToGetApplicationsApi(modelApplications []*model.Application) []*api.Application {
	apiApplications := make([]*api.Application, 0, len(modelApplications))

	for _, modelApplication := range modelApplications {
		apiApplications = append(apiApplications, t.ToApplicationApi(modelApplication))
	}

	return apiApplications
}

func (t *translator) ToApplicationApi(applicationModel *model.Application) *api.Application {
	return &api.Application{
		Name:                  applicationModel.Name,
		Description:           applicationModel.Description,
		Team:                  t.ToApiTeam(applicationModel.Team),
		GitInformation:        t.ToApiGitInformation(applicationModel.GitInformation),
		MonitoringInformation: t.ToMonitoringInformationApi(applicationModel.MonitoringInformation),
		Token:                 t.ToTokenApi(applicationModel.Token),
	}
}

func (t *translator) ToMonitoringInformationApi(monitoringInfo *model.MonitoringInformation) *api.MonitoringInformation {
	if monitoringInfo == nil {
		return nil
	}

	return &api.MonitoringInformation{
		HasOpenAPI:     monitoringInfo.HasOpenApi,
		OpenAPIPath:    monitoringInfo.OpenApiPath,
		HasOpenClient:  monitoringInfo.HasOpenClient,
		OpenClientPath: monitoringInfo.OpenClientPath,
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

func (t *translator) ToApiTeam(teamModel *model.Team) *api.Team {
	if teamModel == nil {
		return nil
	}
	return &api.Team{
		Name:        teamModel.Name,
		Description: teamModel.Description,
	}
}

func (t *translator) ToTokenApi(tokenModel *model.Token) *api.Token {
	if tokenModel == nil {
		return nil
	}

	team := ""
	if tokenModel.Team != nil {
		team = tokenModel.Team.Name
	}

	return &api.Token{
		Name: tokenModel.Name,
		Team: team,
	}
}
