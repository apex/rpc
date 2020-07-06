package goserver

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/schema"
)

// Generate writes the Go server implementations to w.
func Generate(w io.Writer, s *schema.Schema, tracing bool, types string) error {
	// router
	err := writeRouter(w, s, types)
	if err != nil {
		return fmt.Errorf("writing router: %w", err)
	}

	// method stubs
	err = writeMethods(w, s, tracing, types)
	if err != nil {
		return fmt.Errorf("writing methods: %w", err)
	}

	return nil
}

// writeRouter writes the routing implementation to w.
func writeRouter(w io.Writer, s *schema.Schema, types string) error {
	out := fmt.Fprintf
	out(w, "// ServeHTTP implementation.\n")
	out(w, "func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n")
	out(w, "  if r.Method == \"GET\" {\n")
	out(w, "    switch r.URL.Path {\n")
	out(w, "      case \"/_health\":\n")
	out(w, "        rpc.WriteHealth(w, s)\n")
	out(w, "      default:\n")
	out(w, "        rpc.WriteError(w, rpc.BadRequest(\"Invalid method\"))\n")
	out(w, "    }\n")
	out(w, "    return\n")
	out(w, "  }\n\n")
	out(w, "  if r.Method == \"POST\" {\n")
	out(w, "    ctx := rpc.NewRequestContext(r.Context(), r)\n")
	out(w, "    var res interface{}\n")
	out(w, "    var err error\n")
	out(w, "    switch r.URL.Path {\n")
	for _, m := range s.Methods {
		out(w, "      case \"/%s\":\n", m.Name)
		// parse input
		if len(m.Inputs) > 0 {
			out(w, "        var in %s\n", format.GoInputType(types, m.Name))
			out(w, "        err = rpc.ReadRequest(r, &in)\n")
			out(w, "        if err != nil {\n")
			out(w, "          break\n")
			out(w, "        }\n")
			out(w, "        res, err = s.%s(ctx, in)\n", format.JsName(m.Name))
		} else {
			out(w, "        res, err = s.%s(ctx)\n", format.JsName(m.Name))
		}
	}
	out(w, "      default:\n")
	out(w, "        err = rpc.BadRequest(\"Invalid method\")\n")
	out(w, "    }\n")
	out(w, "\n")
	out(w, "    if err != nil {\n")
	out(w, "      rpc.WriteError(w, err)\n")
	out(w, "      return\n")
	out(w, "    }\n")
	out(w, "\n")
	out(w, "    rpc.WriteResponse(w, res)\n")
	out(w, "    return\n")
	out(w, "  }\n")
	out(w, "}\n")
	return nil
}

// writeMethods writes method stubs to w.
func writeMethods(w io.Writer, s *schema.Schema, tracing bool, types string) error {
	out := fmt.Fprintf

	for _, m := range s.Methods {
		out(w, "\n")
		out(w, "// %s %s\n", format.JsName(m.Name), m.Description)

		// method signature
		if len(m.Inputs) > 0 {
			out(w, "func (s *Server) %s(ctx context.Context, in %s) (interface{}, error) {\n", format.JsName(m.Name), format.GoInputType(types, m.Name))
		} else {
			out(w, "func (s *Server) %s(ctx context.Context) (interface{}, error) {\n", format.JsName(m.Name))
		}

		// tracing
		if tracing {
			out(w, "  logs := log.FromContext(ctx).WithField(\"method\", %q)\n\n", m.Name)

			if len(m.Inputs) > 0 {
				out(w, "  logs = logs.WithFields(log.Fields{\n")
				for _, f := range m.Inputs {
					out(w, "    %q: in.%s,\n", f.Name, format.GoName(f.Name))
				}
				out(w, "  })\n")
			}
			out(w, "\n")
		}

		// invoke method
		if len(m.Outputs) > 0 {
			out(w, "  res, err := s.%s", format.GoName(m.Name))
		} else {
			out(w, "  err := s.%s", format.GoName(m.Name))
		}

		if tracing {
			if len(m.Inputs) > 0 {
				out(w, "(log.NewContext(ctx, logs), in)\n")
			} else {
				out(w, "(log.NewContext(ctx, logs))\n")
			}
		} else {
			if len(m.Inputs) > 0 {
				out(w, "(ctx, in)\n")
			} else {
				out(w, "(ctx)\n")
			}
		}

		if len(m.Outputs) > 0 {
			out(w, "  return res, err\n")
		} else {
			out(w, "  return nil, err\n")
		}

		out(w, "}\n")
	}
	out(w, "\n")

	return nil
}
