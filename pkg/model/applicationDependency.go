package model

type ApplicationDependency struct {
	Consumer  *Application
	Provider  *Application
	Reasons   []string
	Endpoints []Endpoint
}

type Endpoint struct {
	Path    string
	Method  string
	Reasons []string
}
