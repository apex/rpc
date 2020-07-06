package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/apex/rpc/schema"
	"github.com/tj/assert"
)

func TestGenerateNoTypes(t *testing.T) {
	exp, err := ioutil.ReadFile("fixtures/todo_no_types_gen.go.fix")
	if err != nil {
		t.Fatalf("error reading no types fixture: %v", err)
	}

	schema, err := schema.Load("../../examples/todo/schema.json")
	if err != nil {
		log.Fatalf("error loading schema: %s", err)
	}

	var act bytes.Buffer
	err = generate(&act, schema, "server", "", true)
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}

	assert.Equal(t, string(exp), act.String())
}

func TestGenerateTypes(t *testing.T) {
	exp, err := ioutil.ReadFile("fixtures/todo_types_gen.go.fix")
	if err != nil {
		t.Fatalf("error reading no types fixture: %v", err)
	}

	schema, err := schema.Load("../../examples/todo/schema.json")
	if err != nil {
		log.Fatalf("error loading schema: %s", err)
	}

	var act bytes.Buffer
	err = generate(&act, schema, "server", "github.com/apex/rpc/cmd/rpc-go-server/api", true)
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}

	assert.Equal(t, string(exp), act.String())
}
