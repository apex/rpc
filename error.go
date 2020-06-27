package rpc

import (
	"net/http"
)

// StatusProvider is the interface used for providing an HTTP status code.
type StatusProvider interface {
	StatusCode() int
}

// TypeProvider is the interface used for providing an error type.
type TypeProvider interface {
	Type() string
}

// ServerError is a server error which implements StatusProvider and TypeProvider.
type ServerError struct {
	status  int
	kind    string
	message string
}

// StatusCode implementation.
func (e ServerError) StatusCode() int {
	return e.status
}

// Type implementation.
func (e ServerError) Type() string {
	return e.kind
}

// Error implementation.
func (e ServerError) Error() string {
	return e.message
}

// Error returns a new ServerError with HTTP status code, kind and message.
func Error(status int, kind, message string) error {
	return ServerError{
		kind:    kind,
		status:  status,
		message: message,
	}
}

// BadRequest returns a new bad request error.
func BadRequest(message string) error {
	return Error(http.StatusBadRequest, "bad_request", message)
}

// Invalid returns a validation error.
func Invalid(message string) error {
	return Error(http.StatusBadRequest, "invalid", message)
}

// serverErrorResponse is an error response.
type serverErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// WriteError writes an error.
//
// If err is a StatusProvider the status code provided
// is used, otherwise it defaults to StatusInternalServerError.
//
// If err is a TypeProvider the type provided is used,
// otherwise it defaults to "internal".
//
// The message in the response uses the Error()
// implementation.
//
func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if e, ok := err.(StatusProvider); ok {
		w.WriteHeader(e.StatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var body serverErrorResponse

	if e, ok := err.(TypeProvider); ok {
		body.Type = e.Type()
	} else {
		body.Type = "internal"
	}

	body.Message = err.Error()
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(body)
}
