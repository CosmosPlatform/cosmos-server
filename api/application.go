package api

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateApplicationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Team        string `json:"team,omitempty"`
}

func (r *CreateApplicationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(0, 500)),
		validation.Field(&r.Team, validation.Length(0, 100)),
	)
}

type CreateApplicationResponse struct {
	Application *Application `json:"application"`
}

type Application struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Team        *Team  `json:"team,omitempty"`
}

type GetApplicationsResponse struct {
	Applications []*Application `json:"applications"`
}
