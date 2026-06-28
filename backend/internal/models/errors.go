package models

// ValidationError represents a domain validation failure with a stable
// machine-readable code and a human-readable message.
type ValidationError struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new ValidationError.
func NewValidationError(code, message string) *ValidationError {
	return &ValidationError{Code: code, Message: message}
}
