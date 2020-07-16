package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/apex/rpc/generators/rubyclient"
	"github.com/apex/rpc/schema"
)

func main() {
	path := flag.String("schema", "schema.json", "Path to the schema file")
	moduleName := flag.String("module", "MyModule", "Name of the module")
	className := flag.String("class", "Client", "Name of the client class")
	flag.Parse()

	s, err := schema.Load(*path)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = generate(os.Stdout, s, *moduleName, *className)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// generate implementation.
func generate(w io.Writer, s *schema.Schema, moduleName, className string) error {
	err := rubyclient.Generate(w, s, moduleName, className)
	if err != nil {
		return fmt.Errorf("generating client: %w", err)
	}

	return nil
}
