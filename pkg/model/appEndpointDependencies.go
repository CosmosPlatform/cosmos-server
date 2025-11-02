package model

type AppEndpointDependencies struct {
	Application *Application
	Endpoints   map[string]bool
}
