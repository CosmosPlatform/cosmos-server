package obj

type Application struct {
	CosmosObj
	Name                string `gorm:"uniqueIndex"`
	Description         string
	TeamID              *int
	Team                *Team `gorm:"foreignKey:TeamID"`
	GitProvider         string
	GitRepositoryOwner  string
	GitRepositoryName   string
	GitRepositoryBranch string
	DependenciesSha     string
	OpenAPISha          string
	HasOpenApi          bool
	OpenApiPath         string
	HasOpenClient       bool
	OpenClientPath      string
}
