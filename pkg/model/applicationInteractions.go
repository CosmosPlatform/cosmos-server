package model

type ApplicationInteractions struct {
	MainApplication       *Application
	ApplicationsToProvide []*Application
	ApplicationsToConsume []*Application
	Interactions          []*ApplicationDependency
}
