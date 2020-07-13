package tsclient_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/apex/rpc/generators/tsclient"
	"github.com/apex/rpc/schema"
	"github.com/tj/assert"
)

func TestGenerate(t *testing.T) {
	exp, err := ioutil.ReadFile("testdata/todo_client.ts")
	assert.NoError(t, err, "reading fixture")

	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = tsclient.Generate(&act, schema, "node-fetch")
	assert.NoError(t, err, "generating")

	// ioutil.WriteFile("testdata/todo_client.ts", act.Bytes(), 0755)

	assert.Equal(t, string(exp), act.String())
}
