package errors

type ErrorType uint

const (
	BadRequest ErrorType = iota
	Unauthorized
	Forbidden
	NotFound
	InternalServerError
	Conflict
)

type ProgramError interface {
	Type() ErrorType
	Error() string
	Details() []string
}

type badRequestError struct {
	error   string
	details []string
}

func NewBadRequestError(error string, details ...string) ProgramError {
	return &badRequestError{
		error:   error,
		details: details,
	}
}

func (e *badRequestError) Type() ErrorType {
	return BadRequest
}

func (e *badRequestError) Error() string {
	return e.error
}

func (e *badRequestError) Details() []string {
	return e.details
}

type unauthorizedError struct {
	error   string
	details []string
}

func NewUnauthorizedError(error string, details ...string) ProgramError {
	return &unauthorizedError{
		error:   error,
		details: details,
	}
}

func (e *unauthorizedError) Type() ErrorType {
	return Unauthorized
}

func (e *unauthorizedError) Error() string {
	return e.error
}

func (e *unauthorizedError) Details() []string {
	return e.details
}

type forbiddenError struct {
	error   string
	details []string
}

func NewForbiddenError(error string, details ...string) ProgramError {
	return &forbiddenError{
		error:   error,
		details: details,
	}
}

func (e *forbiddenError) Type() ErrorType {
	return Forbidden
}

func (e *forbiddenError) Error() string {
	return e.error
}

func (e *forbiddenError) Details() []string {
	return e.details
}

type notFoundError struct {
	error   string
	details []string
}

func NewNotFoundError(error string, details ...string) ProgramError {
	return &notFoundError{
		error:   error,
		details: details,
	}
}

func (e *notFoundError) Type() ErrorType {
	return NotFound
}

func (e *notFoundError) Error() string {
	return e.error
}

func (e *notFoundError) Details() []string {
	return e.details
}

type internalServerError struct {
	error   string
	details []string
}

func NewInternalServerError(error string, details ...string) ProgramError {
	return &internalServerError{
		error:   error,
		details: details,
	}
}

func (e *internalServerError) Type() ErrorType {
	return InternalServerError
}

func (e *internalServerError) Error() string {
	return e.error
}

func (e *internalServerError) Details() []string {
	return e.details
}

type conflictError struct {
	error   string
	details []string
}

func NewConflictError(error string, details ...string) ProgramError {
	return &conflictError{
		error:   error,
		details: details,
	}
}

func (e *conflictError) Type() ErrorType {
	return Conflict
}

func (e *conflictError) Error() string {
	return e.error
}

func (e *conflictError) Details() []string {
	return e.details
}
