package monitoring

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

type Translator interface {
	ToApplicationModel(objApplication *obj.Application) *model.Application
	ToApplicationDependencyObj(modelDependency *model.ApplicationDependency) *obj.ApplicationDependency
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToApplicationModel(applicationObj *obj.Application) *model.Application {
	modelApplication := &model.Application{
		Name:        applicationObj.Name,
		Description: applicationObj.Description,
		Team:        t.ToModelTeam(applicationObj.Team),
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

func (t *translator) ToApplicationDependencyObj(modelDependency *model.ApplicationDependency) *obj.ApplicationDependency {
	if modelDependency == nil {
		return nil
	}

	return &obj.ApplicationDependency{
		Reasons:   modelDependency.Reasons,
		Endpoints: t.toObjEndpoints(modelDependency.Endpoints),
	}
}

func (t *translator) toObjEndpoints(modelEndpoints []model.Endpoint) obj.Endpoints {
	endpoints := make(obj.Endpoints)

	for _, endpoint := range modelEndpoints {
		if _, exists := endpoints[endpoint.Path]; !exists {
			endpoints[endpoint.Path] = make(obj.EndpointMethods)
		}

		endpoints[endpoint.Path][endpoint.Method] = obj.EndpointDetails{
			Reasons: endpoint.Reasons,
		}
	}

	return endpoints
}
