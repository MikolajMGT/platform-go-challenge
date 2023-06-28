package errs

import "errors"

var (
	ValidationError     = errors.New("validation error")
	AuthenticationError = errors.New("failed to authenticate user")
	ProcessingError     = errors.New("processing error")
	AlreadyExistsError  = errors.New("entity already exits")
	CannotBeFoundError  = errors.New("entity cannot be found")
)
