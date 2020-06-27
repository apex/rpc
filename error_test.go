package rpc_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tj/assert"

	"github.com/apex/rpc"
)

// Test error reporting.
func TestWriteError(t *testing.T) {
	t.Run("with a regular error", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteError(w, errors.New("boom"))
		assert.Equal(t, 500, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, "{\n  \"type\": \"internal\",\n  \"message\": \"boom\"\n}", strings.TrimSpace(w.Body.String()))
	})

	t.Run("with a TypeProvider", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteError(w, rpc.Error(400, "invalid_slug", "Invalid team slug"))
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, "{\n  \"type\": \"invalid_slug\",\n  \"message\": \"Invalid team slug\"\n}", strings.TrimSpace(w.Body.String()))
	})

	t.Run("with a StatusProvider", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteError(w, rpc.Error(400, "invalid_slug", "Invalid team slug"))
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Equal(t, 400, w.Code)
	})
}
