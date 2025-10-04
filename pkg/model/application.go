package model

type Application struct {
	Name                  string
	Description           string
	Team                  *Team
	GitInformation        *GitInformation
	MonitoringInformation *MonitoringInformation
}

type GitInformation struct {
	Provider         string
	RepositoryOwner  string
	RepositoryName   string
	RepositoryBranch string
}

type MonitoringInformation struct {
	DependenciesSha string
	OpenAPISha      string
}

type ApplicationUpdate struct {
	Name           *string
	Description    *string
	Team           *string
	GitInformation *GitInformation
}
