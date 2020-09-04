package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tj/assert"

	"github.com/apex/rpc"
)

// DiscardResponseWriter .
type DiscardResponseWriter struct {
}

// Header implementation.
func (d DiscardResponseWriter) Header() http.Header {
	return http.Header{}
}

// Write implementation.
func (d DiscardResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

// WriteHeader implementation.
func (d DiscardResponseWriter) WriteHeader(int) {}

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

// Benchmark responses.
func BenchmarkWriteResponse(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)

	out := struct{ Name, Species, Email string }{
		Name:    "Tobi",
		Species: "Ferret",
		Email:   "tobi@ferret.com",
	}

	var w DiscardResponseWriter
	for i := 0; i < b.N; i++ {
		rpc.WriteResponse(w, out)
	}
}
