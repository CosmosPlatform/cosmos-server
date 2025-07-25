package api

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Team struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type GetTeamsResponse struct {
	Teams []*Team `json:"teams"`
}

type InsertTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (r *InsertTeamRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(0, 500)),
	)
}

type InsertTeamResponse struct {
	Team *Team `json:"team"`
}
