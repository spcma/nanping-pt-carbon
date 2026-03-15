package domain

import "errors"

var (
	ErrMethodologyNotFound      = errors.New("methodology not found")
	ErrMethodologyAlreadyExists = errors.New("methodology already exists")
	ErrMethodologyCodeInvalid   = errors.New("invalid methodology code")
	ErrMethodologyStatusInvalid = errors.New("invalid methodology status")
)
