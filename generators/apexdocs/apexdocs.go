package apexdocs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/internal/schemautil"
	"github.com/apex/rpc/schema"
)

// Generate markdown documentation in the given dir.
func Generate(s *schema.Schema, dir string) error {
	// types
	if err := generateTypes(s, dir); err != nil {
		return fmt.Errorf("generating types: %w", err)
	}

	// methods
	if err := generateMethods(s, dir); err != nil {
		return fmt.Errorf("generating methods: %w", err)
	}

	return nil
}

// generateTypes generates type documentation.
func generateTypes(s *schema.Schema, dir string) error {
	path := filepath.Join(dir, "01-types.md")

	fmt.Printf("  ==> Create %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "---\n")
	fmt.Fprintf(f, "title: Types\n")
	fmt.Fprintf(f, "slug: types\n")
	fmt.Fprintf(f, "teaser: Type documentation for the %s API.\n", s.Name)
	fmt.Fprintf(f, "---\n")

	writeTypes(f, s.TypesSlice())
	return nil
}

// writeTypeIndex writes type index documentation to w.
func writeTypeIndex(w io.Writer, types []schema.Type) {
	for _, t := range types {
		name := format.GoName(t.Name)
		fmt.Fprintf(w, "  - [%s](#%s) — %s\n", name, format.ID(name), t.Description)
	}
}

// writeTypes writes type documentation to w.
func writeTypes(w io.Writer, types []schema.Type) {
	for _, t := range types {
		writeType(w, t)
	}
}

// writeType writes type documentation to w.
func writeType(w io.Writer, t schema.Type) {
	fmt.Fprintf(w, "## %s\n\n", format.GoName(t.Name))
	fmt.Fprintf(w, "The `%s` %s\n\n", format.GoName(t.Name), t.Description)
	writeTableHeader(w, "Name", "Type", "Description")
	for _, f := range t.Properties {
		writeField(w, f)
	}
	writeTypeExamples(w, t.Examples)
}

// writeTypeExamples writes type examples to w.
func writeTypeExamples(w io.Writer, examples []schema.Example) {
	if len(examples) == 0 {
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	fmt.Fprintf(w, "\n### Examples\n\n")
	for _, e := range examples {
		fmt.Fprintf(w, "%s\n\n", e.Description)
		fmt.Fprintf(w, "```json\n")
		enc.Encode(e.Value)
		fmt.Fprintf(w, "```\n\n")
	}
}

// generateMethods generates method documentation.
func generateMethods(s *schema.Schema, dir string) error {
	path := filepath.Join(dir, "02-methods.md")

	fmt.Printf("  ==> Create %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "---\n")
	fmt.Fprintf(f, "title: Methods\n")
	fmt.Fprintf(f, "slug: methods\n")
	fmt.Fprintf(f, "teaser: Method documentation for the %s API.\n", s.Name)
	fmt.Fprintf(f, "---\n")

	writeMethodGroups(f, s)

	return nil
}

// writeMethodGroups writes methods in groups.
func writeMethodGroups(w io.Writer, s *schema.Schema) {
	for _, g := range s.Groups {
		fmt.Fprintf(w, "## %s\n\n", capitalize(g.Name))
		fmt.Fprintf(w, "%s\n\n", g.Description)

		// index
		writeTableHeader(w, "Method", "Description")
		for _, m := range s.Methods {
			if m.Group != g.Name {
				continue
			}
			name := fmt.Sprintf(`[%s](#%s.%s)`, m.Name, format.ID(g.Name), format.ID(m.Name))
			writeTableRow(w, name, capitalize(m.Description))
		}

		// methods
		for _, m := range s.Methods {
			if m.Group != g.Name {
				continue
			}
			writeMethod(w, m)
		}
		fmt.Fprintf(w, "\n")
	}
}

// writeMethod writes method documentation to w.
func writeMethod(w io.Writer, m schema.Method) {
	fmt.Fprintf(w, "### %s\n\n", m.Name)
	fmt.Fprintf(w, "The `%s` method %s\n\n", m.Name, m.Description)

	// inputs
	if len(m.Inputs) > 0 {
		fmt.Fprintf(w, "  Inputs:\n\n")
		writeTableHeader(w, "Name", "Type", "Description")
		for _, f := range m.Inputs {
			writeField(w, f)
		}
		fmt.Fprintf(w, "\n")
	}

	// outputs
	if len(m.Outputs) > 0 {
		fmt.Fprintf(w, "  Outputs:\n\n")
		writeTableHeader(w, "Name", "Type", "Description")
		for _, f := range m.Outputs {
			writeField(w, f)
		}
	}

	writeMethodExamples(w, m.Examples)
	fmt.Fprintf(w, "\n")
}

// writeMethodExamples writes method examples to w.
func writeMethodExamples(w io.Writer, examples []schema.MethodExample) {
	if len(examples) == 0 {
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	if len(examples) > 1 {
		fmt.Fprintf(w, "\n#### Examples\n\n")
	} else {
		fmt.Fprintf(w, "\n#### Example\n\n")
	}

	for _, e := range examples {
		if e.Name != "" {
			fmt.Fprintf(w, "##### %s\n\n", e.Name)
		}

		if e.Description != "" {
			fmt.Fprintf(w, "%s\n\n", e.Description)
		}

		fmt.Fprintf(w, "Input:\n\n")
		if e.Input == nil {
			fmt.Fprintf(w, "  None.\n\n")
		} else {
			fmt.Fprintf(w, "```json\n")
			enc.Encode(e.Input)
			fmt.Fprintf(w, "```\n\n")
		}

		fmt.Fprintf(w, "Output:\n\n")
		if e.Output == nil {
			fmt.Fprintf(w, "  None.\n")
		} else {
			fmt.Fprintf(w, "```json\n")
			enc.Encode(e.Output)
			fmt.Fprintf(w, "```\n\n")
		}
	}
}

// writeField writes a field to w.
func writeField(w io.Writer, f schema.Field) {
	name := fmt.Sprintf("`%s`", f.Name)
	kind := formatType(f.Type)
	if f.Type.Type == "array" {
		kind = fmt.Sprintf("__array__ of %s", formatType(schema.TypeObject(f.Items)))
	}
	writeTableRow(w, name, kind, capitalize(f.Description)+schemautil.FormatEnum(f))
}

// writeTableRow writes a table row to w.
func writeTableRow(w io.Writer, cells ...string) {
	fmt.Fprintf(w, "%s\n", strings.Join(cells, " | "))
}

// writeTableHeader writes a table header to w.
func writeTableHeader(w io.Writer, cells ...string) {
	for i, c := range cells {
		cells[i] = fmt.Sprintf("__%s__", c)
	}
	fmt.Fprintf(w, "%s\n", strings.Join(cells, " | "))
	fmt.Fprintf(w, "%s\n", strings.Repeat("--- | ", len(cells)))
}

// formatType returns a formatted type.
func formatType(t schema.TypeObject) string {
	if t.Ref.Value != "" {
		parts := strings.Split(t.Ref.Value, "/")
		name := format.GoName(parts[len(parts)-1])
		return fmt.Sprintf("[%s](../types/#%s)", name, format.ID(name))
	}

	return fmt.Sprintf("__%s__", t.Type)
}

// capitalize returns a capitalized string.
func capitalize(s string) string {
	return strings.ToUpper(string(s[0])) + string(s[1:])
}
