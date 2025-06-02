package apierrors

import "fmt"

// Error represents a custom API error.
type Error struct {
	Code    string      `json:"code"`              // Application-specific error code
	Message string      `json:"message"`           // Human-readable message
	Details interface{} `json:"details,omitempty"` // Optional field for more specific error details
	cause   error       // The underlying error
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying cause of the error, if any.
func (e *Error) Unwrap() error {
	return e.cause
}

// New creates a new API error.
func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// NewWithDetails creates a new API error with details.
func NewWithDetails(code, message string, details interface{}) *Error {
	return &Error{Code: code, Message: message, Details: details}
}

// Wrap creates a new API error that wraps an existing error.
func Wrap(err error, code, message string) *Error {
	return &Error{Code: code, Message: message, cause: err}
}

// WrapWithDetails creates a new API error that wraps an existing error and includes details.
func WrapWithDetails(err error, code, message string, details interface{}) *Error {
	return &Error{Code: code, Message: message, Details: details, cause: err}
}

// Standard error codes
const (
	CodeInternalError      = "INTERNAL_ERROR"
	CodeValidationError    = "VALIDATION_ERROR"
	CodeNotFound           = "NOT_FOUND"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeConflict           = "CONFLICT"            // e.g., duplicate entry
	CodeUnprocessableEntity = "UNPROCESSABLE_ENTITY" // Business rule violation
	CodeBadRequest         = "BAD_REQUEST"          // General bad request
)

// Common error constructors
func ErrInternal(err error) *Error {
	if err == nil {
		return New(CodeInternalError, "An unexpected internal error occurred.")
	}
	return Wrap(err, CodeInternalError, fmt.Sprintf("An unexpected internal error occurred: %v", err))
}

func ErrValidation(message string, details interface{}) *Error {
	if message == "" {
		message = "Input validation failed."
	}
	return NewWithDetails(CodeValidationError, message, details)
}

func ErrNotFound(resourceName string, resourceID string) *Error {
	return New(CodeNotFound, fmt.Sprintf("%s with ID '%s' not found.", resourceName, resourceID))
}
func ErrNotFoundGeneric(message string) *Error {
	if message == "" {
		message = "The requested resource was not found."
	}
	return New(CodeNotFound, message)
}

func ErrUnauthorized(message string) *Error {
	if message == "" {
		message = "Authentication is required to access this resource."
	}
	return New(CodeUnauthorized, message)
}

func ErrForbidden(message string) *Error {
	if message == "" {
		message = "You do not have permission to perform this action."
	}
	return New(CodeForbidden, message)
}

func ErrConflict(message string) *Error {
	if message == "" {
		message = "A conflict occurred with the current state of the resource."
	}
	return New(CodeConflict, message)
}

func ErrUnprocessableEntity(message string, details interface{}) *Error {
	if message == "" {
		message = "The request was well-formed but was unable to be followed due to semantic errors."
	}
	return NewWithDetails(CodeUnprocessableEntity, message, details)
}

func ErrBadRequest(message string) *Error {
	if message == "" {
		message = "The request could not be understood by the server due to malformed syntax."
	}
	return New(CodeBadRequest, message)
}