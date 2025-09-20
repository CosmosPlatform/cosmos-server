package api

type GetApplicationInteractionsResponse struct {
	MainApplication      string                            `json:"mainApplication"`
	ApplicationsInvolved map[string]ApplicationInformation `json:"applicationsInvolved"`
	Dependencies         []ApplicationDependency           `json:"dependencies"`
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

type GetApplicationsInteractionsResponse struct {
	ApplicationsInvolved map[string]ApplicationInformation `json:"applicationsInvolved"`
	Dependencies         []ApplicationDependency           `json:"dependencies"`
}
