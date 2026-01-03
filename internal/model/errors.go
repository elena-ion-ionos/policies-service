package model

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrIdNotFound          = errors.New("policy id does not exist")
	ErrLimitExceeded       = errors.New("limit exceeded")
	ErrUnprocessableEntity = errors.New("request body is unprocessable")
	ErrFilterInvalid       = errors.New("Accesskey filter must be used or contract number and identity filter")
)
