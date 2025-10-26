package api

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateTokenRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (r *CreateTokenRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
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
