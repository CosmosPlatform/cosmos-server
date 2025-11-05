package storage

import (
	"context"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/storage Service

type Service interface {
	GetUserWithEmail(ctx context.Context, email string) (*obj.User, error)
	GetUserWithRole(ctx context.Context, role string) (*obj.User, error)
	GetUsersWithFilter(ctx context.Context, filter string) ([]*obj.User, error)
	InsertUser(ctx context.Context, user *obj.User) error
	DeleteUser(ctx context.Context, email string) error

	InsertTeam(ctx context.Context, team *obj.Team) error
	GetTeamsWithFilter(ctx context.Context, filter string) ([]*obj.Team, error)
	DeleteTeam(ctx context.Context, name string) error
	GetTeamWithName(ctx context.Context, name string) (*obj.Team, error)
	AddUserToTeam(ctx context.Context, teamName, username string) error
	RemoveUserFromTeam(ctx context.Context, username string) error
	GetTeamMembers(ctx context.Context, teamName string) ([]*obj.User, error)

	InsertApplication(ctx context.Context, application *obj.Application) error
	GetApplicationWithName(ctx context.Context, name string) (*obj.Application, error)
	GetApplicationsByTeam(ctx context.Context, team string) ([]*obj.Application, error)
	GetApplicationsWithFilter(ctx context.Context, filter string) ([]*obj.Application, error)
	DeleteApplicationWithName(ctx context.Context, name string) error
	UpdateApplication(ctx context.Context, application *obj.Application) error

	GetApplicationDependency(ctx context.Context, consumerID, providerID int) (*obj.ApplicationDependency, error)
	GetApplicationDependenciesWithApplicationInvolved(ctx context.Context, applicationName string) ([]*obj.ApplicationDependency, error)
	GetApplicationDependenciesWithFilter(ctx context.Context, filters model.ApplicationDependencyFilter) ([]*obj.ApplicationDependency, error)
	GetApplicationDependenciesByConsumer(ctx context.Context, consumerName string) ([]*obj.ApplicationDependency, error)
	GetApplicationDependenciesByProvider(ctx context.Context, providerName string) ([]*obj.ApplicationDependency, error)
	GetApplicationDependenciesFromGroup(ctx context.Context, group *obj.Group) ([]*obj.ApplicationDependency, error)

	UpsertOpenAPISpecification(ctx context.Context, applicationName string, openAPISpec *obj.ApplicationOpenAPI, applicationOpenApiSHA string) error
	UpdateApplicationDependencies(ctx context.Context, applicationName string, dependenciesToUpsert map[string]*obj.ApplicationDependency, pendingDependencies map[string]*obj.PendingApplicationDependency, dependenciesToDelete []*obj.ApplicationDependency, applicationDependenciesSHA string) error
	CheckPendingDependenciesForApplication(ctx context.Context, applicationName string) error
	GetOpenAPISpecificationByApplicationName(ctx context.Context, applicationName string) (*obj.ApplicationOpenAPI, error)

	GetSentinelSetting(ctx context.Context, name string) (*obj.SentinelSetting, error)
	InsertSentinelSetting(ctx context.Context, setting *obj.SentinelSetting) error
	UpdateSentinelSetting(ctx context.Context, setting *obj.SentinelSetting) error
	GetApplicationsToMonitor(ctx context.Context) ([]*obj.Application, error)

	InsertToken(ctx context.Context, token *obj.Token) error
	GetTokensFromTeam(ctx context.Context, teamName string) ([]*obj.Token, error)
	GetTokenWithNameAndTeamID(ctx context.Context, name string, teamID int) (*obj.Token, error)
	GetAllTokens(ctx context.Context) ([]*obj.Token, error)
	DeleteToken(ctx context.Context, name string, team string) error
	UpdateToken(ctx context.Context, token *obj.Token) error

	CreateGroup(ctx context.Context, name, description string, applications []*obj.Application) error
	GetGroups(ctx context.Context) ([]*obj.Group, error)
	GetGroupByName(ctx context.Context, name string) (*obj.Group, error)
	DeleteGroupByName(ctx context.Context, name string) error
}
