package errors

import (
	"github.com/pkg/errors"
)

// Domain errors
var (
	// ErrInvalidInput represents an error when the input is invalid
	ErrInvalidInput = errors.New("invalid input error")

	// ErrNotFound represents an error when a resource is not found
	ErrNotFound = errors.New("resource not found error")

	// ErrUnauthorized represents an error when a user is not authorized
	ErrUnauthorized = errors.New("unauthorized error")

	// ErrForbidden represents an error when a user is forbidden from accessing a resource
	ErrForbidden = errors.New("forbidden error")

	// ErrInternal represents an internal server error
	ErrInternal = errors.New("internal server error")
)

// External service errors
var (
	// ErrDiscogsAPI represents an error from the Discogs API
	ErrDiscogsAPI = errors.New("discogs API error")

	// ErrSpotifyAPI represents an error from the Spotify API
	ErrSpotifyAPI = errors.New("spotify API error")

	// ErrExternalService represents a generic error from an external service
	ErrExternalService = errors.New("external service error")
)

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// New creates a new error with the given message
func New(message string) error {
	return errors.New(message)
}

// Is reports whether any error in err's chain matches target
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// WithStack annotates err with a stack trace at the point WithStack was called
func WithStack(err error) error {
	return errors.WithStack(err)
}

// Cause returns the underlying cause of the error, if possible
func Cause(err error) error {
	return errors.Cause(err)
}
