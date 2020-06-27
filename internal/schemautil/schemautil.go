package schemautil

import (
	"fmt"
	"strings"

	"github.com/apex/rpc/schema"
)

// ResolveRef returns a resolved reference, or panics.
func ResolveRef(s *schema.Schema, ref schema.Ref) schema.Type {
	name := strings.Replace(ref.Value, "#/types/", "", 1)

	for _, t := range s.Types {
		if t.Name == name {
			return t
		}
	}

	panic(fmt.Sprintf("reference to undefined type %q", ref.Value))
}

// FormatExtra .
func FormatExtra(f schema.Field) string {
	return FormatAttributes(f) + FormatEnum(f)
}

// FormatEnum returns a formatted enum description.
func FormatEnum(f schema.Field) string {
	if f.Enum == nil {
		return ""
	}

	var values []string
	for _, v := range f.Enum {
		values = append(values, `"`+v+`"`)
	}

	return " Must be one of: " + strings.Join(values, ", ") + "."
}

// FormatAttributes returns a formatted field attributes.
func FormatAttributes(f schema.Field) string {
	var attrs []string

	if f.Required {
		attrs = append(attrs, "required")
	}

	if f.ReadOnly {
		attrs = append(attrs, "read-only")
	}

	if len(attrs) == 0 {
		return ""
	}

	return " This field is " + join(attrs, ", ", " and ") + "."
}

// join list.
func join(list []string, delim, separator string) string {
	if len(list) == 0 {
		return ""
	}

	if len(list) == 1 {
		return list[0]
	}

	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	case 2:
		return list[0] + " " + separator + " " + list[1]
	default:
		last := list[len(list)-1]
		list = list[:len(list)-1]
		return strings.Join(list, delim) + delim + separator + " " + last
	}
}
