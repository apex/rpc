package tstypes_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/apex/rpc/generators/tstypes"
	"github.com/apex/rpc/schema"
	"github.com/tj/assert"
)

func TestGenerate(t *testing.T) {
	exp, err := ioutil.ReadFile("testdata/todo_types.ts")
	assert.NoError(t, err, "reading fixture")

	schema, err := schema.Load("../../examples/todo/schema.json")
	assert.NoError(t, err, "loading schema")

	var act bytes.Buffer
	err = tstypes.Generate(&act, schema)
	assert.NoError(t, err, "generating")

	// ioutil.WriteFile("testdata/todo_types.ts", act.Bytes(), 0755)

	assert.Equal(t, string(exp), act.String())
}
