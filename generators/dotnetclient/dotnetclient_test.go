package dotnetclient

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/apex/rpc/schema"
)

func TestGenerate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = Generate(&act, schema, "ApexLogs", "Client")
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_client.cs", act.Bytes())
}
