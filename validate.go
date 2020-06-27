package rpc

import "fmt"

// Validator is the interface used for validating input.
type Validator interface {
	Validate() error
}

// ValidationError is a field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implementation.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s %s", e.Field, e.Message)
}
