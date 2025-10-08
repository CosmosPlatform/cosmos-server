package model

type ApplicationDependency struct {
	Consumer  *Application
	Provider  *Application
	Reasons   []string
	Endpoints Endpoints
}

type Endpoints map[string]EndpointMethods

type EndpointMethods map[string]EndpointDetails

type EndpointDetails struct {
	Reasons []string `json:"reasons,omitempty"`
}

type PendingApplicationDependency struct {
	Consumer     *Application
	ProviderName string
	Reasons      []string
	Endpoints    Endpoints
}
