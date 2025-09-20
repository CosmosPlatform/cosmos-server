package model

type ApplicationInteractions struct {
	MainApplication      string
	ApplicationsInvolved map[string]*Application
	Interactions         []*ApplicationDependency
}

type ApplicationsInteractions struct {
	ApplicationsInvolved map[string]*Application
	Interactions         []*ApplicationDependency
}
