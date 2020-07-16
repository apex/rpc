package phpclient_test

import (
	"bytes"
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go-fixture"

	"github.com/apex/rpc/generators/phpclient"
	"github.com/apex/rpc/schema"
)

func TestGenerate(t *testing.T) {
	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = phpclient.Generate(&act, schema, "Client")
	assert.NoError(t, err, "generating")

	fixture.Assert(t, "todo_client.php", act.Bytes())
}
