package rpc

import (
	"fmt"
	"net/http"
)

// TODO: remove health checking, this can be provided
// in middleware and has nothing to do with RPC :D

// HealthChecker is the interface used for servers providing a health check.
type HealthChecker interface {
	Health() error
}

// WriteHealth responds with 200 OK or invokes the Health() method on the server
// if it implements the HealthChecker interface.
func WriteHealth(w http.ResponseWriter, s interface{}) {
	h, ok := s.(HealthChecker)
	if ok {
		err := h.Health()
		if err != nil {
			http.Error(w, "Health check failed", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "OK")
}
