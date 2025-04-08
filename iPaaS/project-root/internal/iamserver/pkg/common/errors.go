package common

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrAlreadyExists  = errors.New("already exists")

	ErrInvalidArgument = errors.New("invalid argument")
)
