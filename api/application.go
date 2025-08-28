package api

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateApplicationRequest struct {
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Team           string          `json:"team,omitempty"`
	GitInformation *GitInformation `json:"gitInformation,omitempty"`
}

type GitInformation struct {
	Provider         string `json:"provider,omitempty"`
	RepositoryOwner  string `json:"repositoryOwner,omitempty"`
	RepositoryName   string `json:"repositoryName,omitempty"`
	RepositoryBranch string `json:"repositoryBranch,omitempty"`
}

func (r *CreateApplicationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(0, 500)),
		validation.Field(&r.Team, validation.Length(0, 100)),
		validation.Field(&r.GitInformation, validation.When(r.GitInformation != nil,
			validation.Required.Error("git information is required when provided"),
			validation.By(func(value interface{}) error {
				if gi, ok := value.(*GitInformation); ok && gi != nil {
					return validation.ValidateStruct(gi,
						validation.Field(&gi.Provider, validation.Required),
						validation.Field(&gi.RepositoryOwner, validation.Required),
						validation.Field(&gi.RepositoryName, validation.Required),
						validation.Field(&gi.RepositoryBranch, validation.Required),
					)
				}
				return nil
			}),
		)),
	)
}

type CreateApplicationResponse struct {
	Application *Application `json:"application"`
}

type Application struct {
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Team           *Team           `json:"team,omitempty"`
	GitInformation *GitInformation `json:"gitInformation,omitempty"`
}

type GetApplicationResponse struct {
	Application *Application `json:"application"`
}

type GetApplicationsResponse struct {
	Applications []*Application `json:"applications"`
}

type UpdateApplicationRequest struct {
	Name           *string         `json:"name,omitempty"`
	Description    *string         `json:"description,omitempty"`
	Team           *string         `json:"team,omitempty"`
	GitInformation *GitInformation `json:"gitInformation,omitempty"`
}

func (r *UpdateApplicationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.When(r.Name != nil, validation.Length(1, 100))),
		validation.Field(&r.Description, validation.When(r.Description != nil, validation.Length(0, 500))),
		validation.Field(&r.Team, validation.When(r.Team != nil, validation.Length(0, 100))),
		validation.Field(&r.GitInformation, validation.When(r.GitInformation != nil,
			validation.By(func(value interface{}) error {
				if gi, ok := value.(*GitInformation); ok && gi != nil {
					return validation.ValidateStruct(gi,
						validation.Field(&gi.Provider, validation.Required),
						validation.Field(&gi.RepositoryOwner, validation.Required),
						validation.Field(&gi.RepositoryName, validation.Required),
						validation.Field(&gi.RepositoryBranch, validation.Required),
					)
				}
				return nil
			}),
		)),
	)
}

type UpdateApplicationResponse struct {
	Application *Application `json:"application"`
}
