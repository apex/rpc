package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tj/assert"

	"github.com/apex/rpc"
)

// Test responses.
func TestWriteResponse(t *testing.T) {
	t.Run("with no content", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteResponse(w, nil)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, ``, strings.TrimSpace(w.Body.String()))
	})

	t.Run("with no content", func(t *testing.T) {
		w := httptest.NewRecorder()
		rpc.WriteResponse(w, struct {
			Name string `json:"name"`
		}{
			Name: "Tobi",
		})
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{\n  \"name\": \"Tobi\"\n}", strings.TrimSpace(w.Body.String()))
	})
}
