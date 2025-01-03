package core

import "errors"

var (
	// ErrInvalidRequest indicates that the guideline request is invalid
	ErrInvalidRequest = errors.New("invalid guideline request")

	// ErrLanguageNotSupported indicates that the requested programming language is not supported
	ErrLanguageNotSupported = errors.New("programming language not supported")

	// ErrProjectTypeNotSupported indicates that the requested project type is not supported
	ErrProjectTypeNotSupported = errors.New("project type not supported")

	// ErrInternalServer indicates an internal server error occurred
	ErrInternalServer = errors.New("internal server error")
)

// IsNotSupported checks if the error is related to unsupported features
func IsNotSupported(err error) bool {
	return errors.Is(err, ErrLanguageNotSupported) || errors.Is(err, ErrProjectTypeNotSupported)
}

// IsInvalidRequest checks if the error is related to invalid request
func IsInvalidRequest(err error) bool {
	return errors.Is(err, ErrInvalidRequest)
}
