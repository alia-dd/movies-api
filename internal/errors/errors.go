package errors

import "errors"

var (
	ErrNotFound           = errors.New("Record Not Found")
	ErrDuplicateKey       = errors.New("Duplicate Key Violation")
	ErrInvalidInput       = errors.New("Invalid Input")
	ErrInvalidId          = errors.New("Invalid Id")
	ErrorInvalidForceType = errors.New("Force Must Be a Boolean")
	ErrorMarshel          = errors.New("Failed To Marshel Data")
	ErrUnmarshel          = errors.New("Failed To Unmarshel")
	ErrTransaction        = errors.New("Failed To Begin Transaction")
	ErrServerErr          = errors.New("Internal Server Error")
)
