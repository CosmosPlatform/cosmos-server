package storage

import "errors"

var (
	ErrNotFound      = errors.New("storage: item not found")
	ErrAlreadyExists = errors.New("storage: item already exists")
	ErrInvalidInput  = errors.New("storage: invalid input")
	ErrInternal      = errors.New("storage: internal error")
)
