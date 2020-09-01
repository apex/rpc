package tstypes

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/internal/schemautil"
	"github.com/apex/rpc/schema"
)

// Generate writes the TS type implementations to w.
func Generate(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf

	// types
	for _, t := range s.TypesSlice() {
		out(w, "// %s %s\n", format.GoName(t.Name), t.Description)
		out(w, "export interface %s {\n", format.GoName(t.Name))
		writeFields(w, s, t.Properties)
		out(w, "}\n\n")
	}

	// method types
	for _, m := range s.Methods {
		name := format.GoName(m.Name)

		// inputs
		if len(m.Inputs) > 0 {
			out(w, "// %sInput params.\n", name)
			out(w, "interface %sInput {\n", name)
			writeFields(w, s, m.Inputs)
			out(w, "}\n")
		}

		// both
		if len(m.Inputs) > 0 && len(m.Outputs) > 0 {
			out(w, "\n")
		}

		// outputs
		if len(m.Outputs) > 0 {
			out(w, "// %sOutput params.\n", name)
			out(w, "interface %sOutput {\n", name)
			writeFields(w, s, m.Outputs)
			out(w, "}\n")
		}

		out(w, "\n")
	}

	return nil
}

// writeFields to writer.
func writeFields(w io.Writer, s *schema.Schema, fields []schema.Field) {
	for i, f := range fields {
		writeField(w, s, f)
		if i < len(fields)-1 {
			fmt.Fprintf(w, "\n")
		}
	}
}

// writeField to writer.
func writeField(w io.Writer, s *schema.Schema, f schema.Field) {
	fmt.Fprintf(w, "  // %s is %s%s\n", f.Name, f.Description, schemautil.FormatExtra(f))
	if f.Required {
		fmt.Fprintf(w, "  %s: %s\n", f.Name, jsType(s, f))
	} else {
		fmt.Fprintf(w, "  %s?: %s\n", f.Name, jsType(s, f))
	}
}

// jsType returns a JS equivalent type for field f.
func jsType(s *schema.Schema, f schema.Field) string {
	// ref
	if ref := f.Type.Ref.Value; ref != "" {
		t := schemautil.ResolveRef(s, f.Type.Ref)
		return format.GoName(t.Name)
	}

	// type
	switch f.Type.Type {
	case schema.String:
		return "string"
	case schema.Int, schema.Float:
		return "number"
	case schema.Bool:
		return "boolean"
	case schema.Timestamp:
		return "Date"
	case schema.Object:
		return "object"
	case schema.Array:
		return jsType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		}) + "[]"
	default:
		panic("unhandled type")
	}
}
