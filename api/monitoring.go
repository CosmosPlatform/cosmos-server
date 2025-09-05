package api

type GetApplicationInteractionsResponse struct {
	MainApplication       ApplicationInformation   `json:"mainApplication"`
	ApplicationsToProvide []ApplicationInformation `json:"applicationsToProvide"`
	ApplicationsToConsume []ApplicationInformation `json:"applicationsToConsume"`
	Dependencies          []ApplicationDependency  `json:"dependencies"`
}

type ApplicationInformation struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

type ApplicationDependency struct {
	Consumer  string    `json:"consumer"`
	Provider  string    `json:"provider"`
	Reasons   []string  `json:"reasons"`
	Endpoints Endpoints `json:"endpoints"`
}

type Endpoints map[string]EndpointMethods

type EndpointMethods map[string]EndpointDetails

type EndpointDetails struct {
	Reasons []string `json:"reasons"`
}
