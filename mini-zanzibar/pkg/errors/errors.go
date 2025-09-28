package errors

import "fmt"

// CustomError represents a custom error with additional context
type CustomError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common error codes
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeDatabaseError = "DATABASE_ERROR"
)

// NewValidationError creates a new validation error
func NewValidationError(message string, details map[string]interface{}) *CustomError {
	return &CustomError{
		Code:    ErrCodeValidation,
		Message: message,
		Details: details,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) *CustomError {
	return &CustomError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Details: map[string]interface{}{"resource": resource},
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *CustomError {
	return &CustomError{
		Code:    ErrCodeUnauthorized,
		Message: message,
		Details: nil,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *CustomError {
	return &CustomError{
		Code:    ErrCodeForbidden,
		Message: message,
		Details: nil,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string) *CustomError {
	return &CustomError{
		Code:    ErrCodeInternal,
		Message: message,
		Details: nil,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation string, err error) *CustomError {
	return &CustomError{
		Code:    ErrCodeDatabaseError,
		Message: fmt.Sprintf("Database operation failed: %s", operation),
		Details: map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		},
	}
}
