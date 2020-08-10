package gotypes

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/internal/schemautil"
	"github.com/apex/rpc/schema"
)

var utils = `// oneOf returns true if s is in the values.
func oneOf(s string, values []string) bool {
  for _, v := range values {
		if s == v {
			return true
		}
	}
	return false
}`

// Generate writes the Go type implementations to w, with optional validation methods.
func Generate(w io.Writer, s *schema.Schema, validate bool) error {
	out := fmt.Fprintf

	// default tags
	if s.Go.Tags == nil {
		s.Go.Tags = []string{"json"}
	}

	// types
	for _, t := range s.TypesSlice() {
		out(w, "// %s %s\n", format.GoName(t.Name), t.Description)
		out(w, "type %s struct {\n", format.GoName(t.Name))
		writeFields(w, s, t.Properties)
		out(w, "}\n\n")
		if validate {
			writeValidation(w, format.GoName(t.Name), t.Properties)
			out(w, "\n")
		}
	}

	// methods
	for _, m := range s.Methods {
		name := format.GoName(m.Name)

		// inputs
		if len(m.Inputs) > 0 {
			out(w, "// %sInput params.\n", name)
			out(w, "type %sInput struct {\n", name)
			writeFields(w, s, m.Inputs)
			out(w, "}\n")
			if validate {
				out(w, "\n")
				writeValidation(w, name+"Input", m.Inputs)
			}
		}

		// both
		if len(m.Inputs) > 0 && len(m.Outputs) > 0 {
			out(w, "\n")
		}

		// outputs
		if len(m.Outputs) > 0 {
			out(w, "// %sOutput params.\n", name)
			out(w, "type %sOutput struct {\n", name)
			writeFields(w, s, m.Outputs)
			out(w, "}\n")
		}

		out(w, "\n")
	}

	if validate {
		out(w, "\n%s\n", utils)
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
	fmt.Fprintf(w, "  // %s is %s%s\n", format.GoName(f.Name), f.Description, schemautil.FormatExtra(f))
	fmt.Fprintf(w, "  %s %s %s\n", format.GoName(f.Name), goType(s, f), fieldTags(f, s.Go.Tags))
}

// goType returns a Go equivalent type for field f.
func goType(s *schema.Schema, f schema.Field) string {
	// ref
	if ref := f.Type.Ref.Value; ref != "" {
		t := schemautil.ResolveRef(s, f.Type.Ref)
		return format.GoName(t.Name)
	}

	// type
	switch f.Type.Type {
	case schema.String:
		return "string"
	case schema.Int:
		return "int"
	case schema.Bool:
		return "bool"
	case schema.Float:
		return "float64"
	case schema.Timestamp:
		return "time.Time"
	case schema.Object:
		return "map[string]interface{}"
	case schema.Array:
		return "[]" + goType(s, schema.Field{
			Type: schema.TypeObject(f.Items),
		})
	default:
		panic("unhandled type")
	}
}

// fieldTags returns tags for a field.
func fieldTags(f schema.Field, tags []string) string {
	var pairs [][]string

	for _, tag := range tags {
		pairs = append(pairs, []string{tag, f.Name})
	}

	return formatTags(pairs)
}

// formatTags returns field tags.
func formatTags(tags [][]string) string {
	var s []string
	for _, t := range tags {
		if len(t) == 2 {
			s = append(s, fmt.Sprintf("%s:%q", t[0], t[1]))
		} else {
			s = append(s, fmt.Sprintf("%s", t[0]))
		}
	}
	return fmt.Sprintf("`%s`", strings.Join(s, " "))
}

// writeValidation writes a validation method implementation to w.
func writeValidation(w io.Writer, name string, fields []schema.Field) error {
	out := fmt.Fprintf
	recv := strings.ToLower(name)[0]
	out(w, "// Validate implementation.\n")
	out(w, "func (%c *%s) Validate() error {\n", recv, name)
	for _, f := range fields {
		writeFieldDefaults(w, f, recv)
		writeFieldValidation(w, f, recv)
	}
	out(w, "  return nil\n")
	out(w, "}\n")
	return nil
}

// writeFieldDefaults writes field defaults to w.
func writeFieldDefaults(w io.Writer, f schema.Field, recv byte) error {
	// TODO: write out a separate Default() method?
	if f.Default == nil {
		return nil
	}

	out := fmt.Fprintf
	name := format.GoName(f.Name)

	switch f.Type.Type {
	case schema.Int:
		out(w, "  if %c.%s == 0 {\n", recv, name)
		out(w, "    %c.%s = %v\n", recv, name, f.Default)
		out(w, "  }\n\n")
	case schema.String:
		out(w, "  if %c.%s == \"\" {\n", recv, name)
		out(w, "    %c.%s = %q\n", recv, name, f.Default)
		out(w, "  }\n\n")
	}

	return nil
}

// writeFieldValidation writes field validation to w.
func writeFieldValidation(w io.Writer, f schema.Field, recv byte) error {
	out := fmt.Fprintf
	name := format.GoName(f.Name)

	writeError := func(msg string) {
		out(w, "    return rpc.ValidationError{ Field: %q, Message: %q }\n", f.Name, msg)
	}

	// required
	if f.Required {
		switch f.Type.Type {
		case schema.Int:
			out(w, "  if %c.%s == 0 {\n", recv, name)
			writeError("is required")
			out(w, "  }\n\n")
		case schema.String:
			out(w, "  if %c.%s == \"\" {\n", recv, name)
			writeError("is required")
			out(w, "  }\n\n")
		case schema.Array, schema.Object:
			out(w, "  if %c.%s == nil {\n", recv, name)
			writeError("is required")
			out(w, "  }\n\n")
		case schema.Timestamp:
			out(w, "  if %c.%s.IsZero() {\n", recv, name)
			writeError("is required")
			out(w, "  }\n\n")
		}
	}

	// enums
	if f.Type.Type == schema.String && f.Enum != nil {
		field := fmt.Sprintf("%c.%s", recv, name)
		out(w, "  if %s != \"\" && !oneOf(%s, %s) {\n", field, field, formatSlice(f.Enum))
		writeError(fmt.Sprintf("must be one of: %s", formatEnum(f.Enum)))
		out(w, "  }\n\n")
	}

	// validate the children of non-primitive arrays
	// TODO: HasRef() or similar?
	if f.Type.Type == schema.Array && f.Items.Ref.Value != "" {
		out(w, "  for i, v := range %c.%s {\n", recv, name)
		out(w, "    if err := v.Validate(); err != nil {\n")
		out(w, "      return fmt.Errorf(\"element %%d: %%s\", i, err.Error())\n")
		out(w, "    }\n")
		out(w, "  }\n\n")
	}

	return nil
}

// formatSlice returns a formatted slice from enum.
func formatSlice(values []string) string {
	var vals []string
	for _, l := range values {
		vals = append(vals, strconv.Quote(l))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(vals, ", "))
}

// formatEnum returns a formatted enum values.
func formatEnum(values []string) string {
	var vals []string
	for _, l := range values {
		vals = append(vals, strconv.Quote(l))
	}
	return strings.Join(vals, ", ")
}
