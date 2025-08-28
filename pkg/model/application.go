package model

type Application struct {
	Name           string
	Description    string
	Team           *Team
	GitInformation *GitInformation
}

type GitInformation struct {
	Provider         string
	RepositoryOwner  string
	RepositoryName   string
	RepositoryBranch string
}
