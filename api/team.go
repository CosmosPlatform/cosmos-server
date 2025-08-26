package api

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Team struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type GetTeamsResponse struct {
	Teams []*Team `json:"teams"`
}

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (r *CreateTeamRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(0, 500)),
	)
}

type CreateTeamResponse struct {
	Team *Team `json:"team"`
}

type AddUserToTeamRequest struct {
	Email string `json:"email"`
}

func (r *AddUserToTeamRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, validation.Length(1, 100), is.EmailFormat),
	)
}

type GetTeamResponse struct {
	Team *Team `json:"team"`
}
