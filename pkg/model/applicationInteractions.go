package model

type ApplicationsInteractions struct {
	ApplicationsInvolved map[string]*Application
	Interactions         []*ApplicationDependency
}
