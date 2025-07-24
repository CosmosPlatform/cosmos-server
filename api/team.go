package api

type Team struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type GetTeamsResponse struct {
	Teams []*Team `json:"teams"`
}
