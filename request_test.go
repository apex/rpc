package rpc_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tj/assert"

	"github.com/apex/rpc"
)

// Test requests.
func TestReadRequest(t *testing.T) {
	t.Run("with a no content-type", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", strings.NewReader(`{ "name": "Tobi" }`))
		var in struct{ Name string }
		err := rpc.ReadRequest(r, &in)
		assert.EqualError(t, err, `Unsupported request Content-Type, must be application/json`)
	})

	t.Run("with malformed JSON", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", strings.NewReader(`{ "name": "Tobi`))
		r.Header.Set("Content-Type", "application/json")
		var in struct{ Name string }
		err := rpc.ReadRequest(r, &in)
		assert.EqualError(t, err, `Failed to parse malformed request body, must be a valid JSON object`)
	})

	t.Run("with JSON array", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", strings.NewReader(`[{}]`))
		r.Header.Set("Content-Type", "application/json")
		var in struct{ Name string }
		err := rpc.ReadRequest(r, &in)
		assert.EqualError(t, err, `Failed to parse malformed request body, must be a valid JSON object`)
	})

	t.Run("with a json body", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", strings.NewReader(`{ "name": "Tobi" }`))
		r.Header.Set("Content-Type", "application/json")
		var in struct{ Name string }
		err := rpc.ReadRequest(r, &in)
		assert.NoError(t, err, "parsing")
		assert.Equal(t, "Tobi", in.Name)
	})
}

// Benchmark requests.
func BenchmarkReadRequest(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest("GET", "/", strings.NewReader(`{ "name": "Tobi", "species": "ferret", "email": "tobi@apex.sh" }`))
		r.Header.Set("Content-Type", "application/json")
		var in struct{ Name, Species, Email string }
		err := rpc.ReadRequest(r, &in)
		if err != nil {
			b.Fatal(err)
		}
	}
}
