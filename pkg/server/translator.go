package server

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	errorUtilities "errors"
)

type Translator interface {
	ToApiError(error error) api.ErrorResponse
	ToStatusCode(programError errors.ProgramError) int
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToApiError(receivedError error) api.ErrorResponse {
	var programError errors.ProgramError
	if errorUtilities.As(receivedError, &programError) {
		return api.ErrorResponse{
			Error:      programError.Error(),
			Details:    programError.Details(),
			StatusCode: t.ToStatusCode(programError),
		}
	}

	return api.ErrorResponse{
		Error:      receivedError.Error(),
		Details:    []string{"An unexpected error occurred."},
		StatusCode: 500,
	}
}

func (t *translator) ToStatusCode(programError errors.ProgramError) int {
	switch programError.Type() {
	case errors.BadRequest:
		return 400
	case errors.Unauthorized:
		return 401
	case errors.Forbidden:
		return 403
	case errors.NotFound:
		return 404
	case errors.InternalServerError:
		return 500
	default:
		return 500
	}
}
