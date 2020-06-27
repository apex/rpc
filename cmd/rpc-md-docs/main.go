package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/apex/rpc/generators/mddocs"
	"github.com/apex/rpc/schema"
)

func main() {
	path := flag.String("schema", "schema.json", "Path to the schema file")
	out := flag.String("output", "docs", "Output directory")
	flag.Parse()

	s, err := schema.Load(*path)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	println()
	defer println()

	err = mddocs.Generate(s, *out)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	fmt.Printf("  ==> Complete\n")
}
