package api

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var tokenNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

type CreateTokenRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (r *CreateTokenRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required,
			validation.Match(applicationNameRegex).Error("name can only contain letters, numbers, and hyphens"),
		),
		validation.Field(&r.Value, validation.Required),
	)
}

type GetTokensResponse struct {
	Tokens []*Token `json:"tokens"`
}

type Token struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

type UpdateTokenRequest struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

func (r *UpdateTokenRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.When(r.Name != nil,
				validation.Required,
				validation.Match(applicationNameRegex).Error("name can only contain letters, numbers, and hyphens"),
			),
		),
		validation.Field(&r.Value,
			validation.When(r.Value != nil,
				validation.Required,
			),
		),
	)
}
