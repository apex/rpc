package tsclient_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/apex/rpc/generators/tsclient"
	"github.com/apex/rpc/schema"
)

func TestGenerate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = tsclient.Generate(&act, schema, "node-fetch")
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_client.ts", act.Bytes())
}
