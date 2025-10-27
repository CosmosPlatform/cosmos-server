package api

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var applicationNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

type CreateApplicationRequest struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	Team                  string                 `json:"team,omitempty"`
	GitInformation        *GitInformation        `json:"gitInformation,omitempty"`
	MonitoringInformation *MonitoringInformation `json:"monitoringInformation,omitempty"`
	TokenName             string                 `json:"tokenName,omitempty"`
}

type GitInformation struct {
	Provider         string `json:"provider,omitempty"`
	RepositoryOwner  string `json:"repositoryOwner,omitempty"`
	RepositoryName   string `json:"repositoryName,omitempty"`
	RepositoryBranch string `json:"repositoryBranch,omitempty"`
}

type MonitoringInformation struct {
	HasOpenAPI     bool   `json:"hasOpenAPI"`
	OpenAPIPath    string `json:"openAPIPath"`
	HasOpenClient  bool   `json:"hasOpenClient"`
	OpenClientPath string `json:"openClientPath"`
}

func (r *CreateApplicationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required,
			validation.Length(1, 100),
			validation.Match(applicationNameRegex).Error("name can only contain letters, numbers, and hyphens"),
		),
		validation.Field(&r.Description, validation.Length(0, 500)),
		validation.Field(&r.Team, validation.Length(0, 100)),
		validation.Field(&r.GitInformation, validation.When(r.MonitoringInformation != nil, validation.Required.Error("git information is required when monitoring is provided"))),
		validation.Field(&r.GitInformation, validation.When(r.GitInformation != nil,
			validation.Required.Error("git information is required when provided"),
			validation.By(func(value any) error {
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
		validation.Field(&r.MonitoringInformation, validation.When(r.MonitoringInformation != nil,
			validation.Required.Error("monitoring information is required when provided"),
			validation.By(func(value any) error {
				if mi, ok := value.(*MonitoringInformation); ok && mi != nil {
					return validation.ValidateStruct(mi,
						validation.Field(&mi.OpenAPIPath, validation.When(mi.HasOpenAPI, validation.Required)),
						validation.Field(&mi.OpenClientPath, validation.When(mi.HasOpenClient, validation.Required)),
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
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	Team                  *Team                  `json:"team,omitempty"`
	GitInformation        *GitInformation        `json:"gitInformation,omitempty"`
	MonitoringInformation *MonitoringInformation `json:"monitoringInformation,omitempty"`
	Token                 *Token                 `json:"token,omitempty"`
}

type GetApplicationResponse struct {
	Application *Application `json:"application"`
}

type GetApplicationsResponse struct {
	Applications []*Application `json:"applications"`
}

type UpdateApplicationRequest struct {
	Name                  *string                `json:"name,omitempty"`
	Description           *string                `json:"description,omitempty"`
	Team                  *string                `json:"team,omitempty"`
	GitInformation        *GitInformation        `json:"gitInformation,omitempty"`
	MonitoringInformation *MonitoringInformation `json:"monitoringInformation,omitempty"`
}

func (r *UpdateApplicationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.When(r.Name != nil, validation.Length(1, 100))),
		validation.Field(&r.Description, validation.When(r.Description != nil, validation.Length(0, 500))),
		validation.Field(&r.Team, validation.When(r.Team != nil, validation.Length(0, 100))),
		validation.Field(&r.GitInformation, validation.When(r.MonitoringInformation != nil, validation.Required.Error("git information is required when monitoring is provided"))),
		validation.Field(&r.GitInformation, validation.When(r.GitInformation != nil,
			validation.By(func(value any) error {
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
		validation.Field(&r.MonitoringInformation, validation.When(r.MonitoringInformation != nil,
			validation.Required.Error("monitoring information is required when provided"),
			validation.By(func(value any) error {
				if mi, ok := value.(*MonitoringInformation); ok && mi != nil {
					return validation.ValidateStruct(mi,
						validation.Field(&mi.OpenAPIPath, validation.When(mi.HasOpenAPI, validation.Required)),
						validation.Field(&mi.OpenClientPath, validation.When(mi.HasOpenClient, validation.Required)),
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
