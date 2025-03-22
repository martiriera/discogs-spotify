package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidInput = New("invalid input error")

	ErrNotFound = New("resource not found error")

	ErrUnauthorized = New("unauthorized error")

	ErrForbidden = New("forbidden error")

	ErrInternal = New("internal server error")
)

func New(message string) error {
	return errors.New(message)
}

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Is reports whether any error in err's chain matches target
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Cause returns the underlying cause of the error, if possible
func Cause(err error) error {
	return errors.Cause(err)
}
