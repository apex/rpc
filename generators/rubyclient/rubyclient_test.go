package rubyclient_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/apex/rpc/generators/rubyclient"
	"github.com/apex/rpc/schema"
)

func TestGenerate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = rubyclient.Generate(&act, schema, "Todo", "Client")
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_client.rb", act.Bytes())
}
