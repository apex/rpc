//go:generate file2go -in schema.json -pkg schema

// Package schema provides the Apex RPC schema.
package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/xeipuuv/gojsonschema"
)

// ValidationError is a validation error.
type ValidationError struct {
	Result *gojsonschema.Result
}

// Error implementation.
func (e ValidationError) Error() (s string) {
	s = "validation failed:\n"
	for _, e := range e.Result.Errors() {
		s += fmt.Sprintf("  - %s\n", e)
	}
	return
}

// Kind is a value type.
type Kind string

// Types available.
const (
	String    Kind = "string"
	Bool      Kind = "boolean"
	Int       Kind = "integer"
	Float     Kind = "float"
	Array     Kind = "array"
	Object    Kind = "object"
	Timestamp Kind = "timestamp"
)

// Ref model.
type Ref struct {
	Value string `json:"$ref"`
}

// TypeObject model.
type TypeObject struct {
	Type Kind `json:"type"`
	Ref
}

// UnmarshalJSON implementation.
func (t *TypeObject) UnmarshalJSON(b []byte) error {
	// "type": { "$ref": ... }
	if b[0] == '{' {
		var r Ref
		if err := json.Unmarshal(b, &r); err != nil {
			return err
		}
		t.Ref = r
		return nil
	}

	// "type": "integer"
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t.Type = Kind(s)
	return nil
}

// ItemsObject model.
type ItemsObject struct {
	Type Kind `json:"type"`
	Ref
}

// Schema model.
type Schema struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Methods     []Method        `json:"methods"`
	Groups      []Group         `json:"groups"`
	Types       map[string]Type `json:"types"`
	Go          struct {
		Tags []string `json:"tags"`
	} `json:"go"`
}

// Method model.
type Method struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Private     bool            `json:"private"`
	Group       string          `json:"group"`
	Inputs      []Field         `json:"inputs"`
	Outputs     []Field         `json:"outputs"`
	Examples    []MethodExample `json:"examples"`
}

// MethodExample model.
type MethodExample struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Input       interface{} `json:"input"`
	Output      interface{} `json:"output"`
}

// Field model.
type Field struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	ReadOnly    bool        `json:"readonly"`
	Default     interface{} `json:"default"`
	Type        TypeObject  `json:"type"`
	Items       ItemsObject `json:"items"`
	Enum        []string    `json:"enum"`
}

// Type model.
type Type struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	Properties  []Field   `json:"properties"`
	Examples    []Example `json:"examples"`
}

// Example model.
type Example struct {
	Description string      `json:"description"`
	Value       interface{} `json:"value"`
}

// Group model.
type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

// TypesSlice returns a sorted slice of types.
func (s Schema) TypesSlice() (v []Type) {
	for _, t := range s.Types {
		// sort fields
		sort.Slice(t.Properties, func(i, j int) bool {
			a := t.Properties[i]
			b := t.Properties[j]
			return a.Name < b.Name
		})

		v = append(v, t)
	}

	// sort types
	sort.Slice(v, func(i, j int) bool {
		return v[i].Name < v[j].Name
	})

	return
}

// Load returns a schema loaded and validated from path.
func Load(path string) (*Schema, error) {
	// TODO: bake into the binary with Go's native 'embed' stuff once it's available
	schema := gojsonschema.NewBytesLoader(SchemaJson)
	doc := gojsonschema.NewReferenceLoader("file://" + path)

	// validate
	result, err := gojsonschema.Validate(schema, doc)
	if err != nil {
		return nil, err
	}

	if !result.Valid() {
		return nil, &ValidationError{
			Result: result,
		}
	}

	var s Schema

	// open
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// unmarshal
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return nil, err
	}

	// populate type names
	for k, v := range s.Types {
		v.Name = k
		s.Types[k] = v
	}

	// sort groups
	sort.Slice(s.Groups, func(i, j int) bool {
		a := s.Groups[i]
		b := s.Groups[j]
		return a.Name < b.Name
	})

	// sort methods
	sort.Slice(s.Methods, func(i, j int) bool {
		a := s.Methods[i]
		b := s.Methods[j]
		return a.Name < b.Name
	})

	// sort method inputs & outputs
	for _, m := range s.Methods {
		sort.Slice(m.Inputs, func(i, j int) bool {
			a := m.Inputs[i]
			b := m.Inputs[j]
			return a.Name < b.Name
		})

		sort.Slice(m.Outputs, func(i, j int) bool {
			a := m.Outputs[i]
			b := m.Outputs[j]
			return a.Name < b.Name
		})
	}

	return &s, nil
}

// IsBuiltin returns true if the type is built-in.
func IsBuiltin(kind Kind) bool {
	switch kind {
	case String, Int, Bool, Float, Array, Object, Timestamp:
		return true
	default:
		return false
	}
}
