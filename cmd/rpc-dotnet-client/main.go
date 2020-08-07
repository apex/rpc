package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/apex/rpc/generators/dotnetclient"
	"github.com/apex/rpc/schema"
)

func main() {
	path := flag.String("schema", "schema.json", "Path to the schema file")
	namespaceName := flag.String("namespace", "MyNamespace", "Name of the namespace")
	className := flag.String("class", "Client", "Name of the client class")
	flag.Parse()

	s, err := schema.Load(*path)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = generate(os.Stdout, s, *namespaceName, *className)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// generate implementation.
func generate(w io.Writer, s *schema.Schema, namespace, className string) error {
	err := dotnetclient.Generate(w, s, namespace, className)
	if err != nil {
		return fmt.Errorf("generating client: %w", err)
	}

	return nil
}
