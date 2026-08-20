package errors

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrDuplicateKey = errors.New("duplicate key violation")
	ErrInvalidInput = errors.New("invalid input")
)
