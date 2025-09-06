package model

type ApplicationInteractions struct {
	MainApplication      string
	ApplicationsInvolved map[string]*Application
	Interactions         []*ApplicationDependency
}
