package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/apex/rpc/schema"
	"github.com/tj/assert"
)

func TestGenerateNoTypes(t *testing.T) {
	exp, err := ioutil.ReadFile("fixtures/todo_no_types_gen.go.fix")
	assert.NoError(t, err, "reading fixture")

	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = generate(&act, schema, "server", "", true)
	assert.NoError(t, err, "generating")

	assert.Equal(t, string(exp), act.String())
}

func TestGenerateTypes(t *testing.T) {
	exp, err := ioutil.ReadFile("fixtures/todo_types_gen.go.fix")
	assert.NoError(t, err, "reading fixture")

	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = generate(&act, schema, "server", "github.com/apex/rpc/cmd/rpc-go-server/api", true)
	assert.NoError(t, err, "generating")

	assert.Equal(t, string(exp), act.String())
}
