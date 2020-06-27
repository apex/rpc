package rpc_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/tj/assert"
	
	"github.com/apex/rpc"
)

// healthChecker implementation.
type healthChecker struct {
	err error
}

// Health implementation.
func (h healthChecker) Health() error {
	return h.err
}

// Test health checks.
func TestWriteHealth(t *testing.T) {
	t.Run("without a HealthChecker", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteHealth(w, nil)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK\n", w.Body.String())
	})

	t.Run("with a HealthChecker passing", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteHealth(w, healthChecker{})
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "OK\n", w.Body.String())
	})

	t.Run("with a HealthChecker failing", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteHealth(w, healthChecker{errors.New("boom")})
		assert.Equal(t, 500, w.Code)
		assert.Equal(t, "Health check failed\n", w.Body.String())
	})
}
