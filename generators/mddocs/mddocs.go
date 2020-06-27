package mddocs

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
	// types dir
	typesDir := filepath.Join(dir, "types")
	if err := os.MkdirAll(typesDir, 0755); err != nil {
		return err
	}

	// types index
	if err := generateTypesIndex(s.TypesSlice(), typesDir); err != nil {
		return fmt.Errorf("generating types index: %w", err)
	}

	// types
	for _, t := range s.Types {
		if err := generateType(t, typesDir); err != nil {
			return fmt.Errorf("generating type: %w", err)
		}
	}

	// methods dir
	methodsDir := filepath.Join(dir, "methods")
	if err := os.MkdirAll(methodsDir, 0755); err != nil {
		return err
	}

	// methods index
	if err := generateMethodsIndex(s, methodsDir); err != nil {
		return fmt.Errorf("generating methods index: %w", err)
	}

	// methods
	for _, m := range s.Methods {
		if err := generateMethod(m, methodsDir); err != nil {
			return fmt.Errorf("generating methods: %w", err)
		}
	}

	return nil
}

// generateTypesIndex generates type index documentation.
func generateTypesIndex(types []schema.Type, dir string) error {
	path := filepath.Join(dir, "index.md")

	fmt.Printf("  ==> Create %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writeTypeIndex(f, types)
	return nil
}

// writeTypeIndex writes type index documentation to w.
func writeTypeIndex(w io.Writer, types []schema.Type) {
	fmt.Fprintf(w, "# Types\n\n")
	for _, t := range types {
		name := format.GoName(t.Name)
		fmt.Fprintf(w, "  - [%s](./%s.md) — %s\n", name, name, t.Description)
	}
}

// generateType generates type documentation.
func generateType(t schema.Type, dir string) error {
	path := filepath.Join(dir, format.GoName(t.Name)+".md")

	fmt.Printf("  ==> Create %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writeType(f, t)
	return nil
}

// writeType writes type documentation to w.
func writeType(w io.Writer, t schema.Type) {
	fmt.Fprintf(w, "# %s\n\n", format.GoName(t.Name))
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

	fmt.Fprintf(w, "\n## Examples\n\n")
	for _, e := range examples {
		fmt.Fprintf(w, "%s\n\n", e.Description)
		fmt.Fprintf(w, "```json\n")
		enc.Encode(e.Value)
		fmt.Fprintf(w, "```\n\n")
	}
}

// generateMethodsIndex generates method index documentation.
func generateMethodsIndex(s *schema.Schema, dir string) error {
	path := filepath.Join(dir, "index.md")

	fmt.Printf("  ==> Create %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writeMethodIndex(f, s)
	return nil
}

// writeMethodIndex writes method index documentation to w.
func writeMethodIndex(w io.Writer, s *schema.Schema) {
	fmt.Fprintf(w, "# Methods\n\n")
	for _, g := range s.Groups {
		fmt.Fprintf(w, "## %s\n\n", capitalize(g.Name))
		fmt.Fprintf(w, "%s\n\n", g.Description)
		for _, m := range s.Methods {
			if m.Group != g.Name {
				continue
			}
			fmt.Fprintf(w, "  - [%s](./%s.md) — %s\n", m.Name, m.Name, m.Description)
		}
		fmt.Fprintf(w, "\n")
	}
}

// generateMethod generates method documentation.
func generateMethod(m schema.Method, dir string) error {
	path := filepath.Join(dir, m.Name+".md")

	fmt.Printf("  ==> Create %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writeMethod(f, m)
	return nil
}

// writeMethod writes method documentation to w.
func writeMethod(w io.Writer, m schema.Method) {
	fmt.Fprintf(w, "# %s\n\n", m.Name)
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
		fmt.Fprintf(w, "\n## Examples\n\n")
	} else {
		fmt.Fprintf(w, "\n## Example\n\n")
	}

	for _, e := range examples {
		if e.Name != "" {
			fmt.Fprintf(w, "### %s\n\n", e.Name)
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
	writeTableRow(w, name, kind, capitalize(f.Description)+schemautil.FormatExtra(f))
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
		return fmt.Sprintf("[%s](../types/%s.md)", name, name)
	}

	return fmt.Sprintf("__%s__", t.Type)
}

// capitalize returns a capitalized string.
func capitalize(s string) string {
	return strings.ToUpper(string(s[0])) + string(s[1:])
}
