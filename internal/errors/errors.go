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
	ErrTransactionStart   = errors.New("Failed To Begin Transaction")
	ErrTransactionCommit  = errors.New("Failed To Commit Transaction")
	ErrServerErr          = errors.New("Internal Server Error")
	ErrWrongApiKey        = errors.New("Invalid Api Key ")
	// ErrMissingPageNumver  = errors.New("Page Numvber not provided Data Missing")
	// ErrMissingLimit       = errors.New("Data Limit no Data Missing")
)
