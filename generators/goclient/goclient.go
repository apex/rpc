package goclient

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/schema"
)

var call = `// Error is an error returned by the client.
type Error struct {
	Status     string
	StatusCode int
	Type       string
	Message    string
}

// Error implementation.
func (e Error) Error() string {
	if e.Type == "" {
		return fmt.Sprintf("%s: %d", e.Status, e.StatusCode)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// call implementation.
func call(client *http.Client, authToken, endpoint, method string, in, out interface{}) error {
	var body io.Reader

	// default client
	if client == nil {
		client = http.DefaultClient
	}

	// input params
	if in != nil {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(in)
		if err != nil {
			return fmt.Errorf("encoding: %w", err)
		}
		body = &buf
	}

	// POST request
	req, err := http.NewRequest("POST", endpoint+"/"+method, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// auth token
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	// response
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// error
	if res.StatusCode >= 300 {
		var e Error
		if res.Header.Get("Content-Type") == "application/json" {
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				return err
			}
		}
		e.Status = http.StatusText(res.StatusCode)
		e.StatusCode = res.StatusCode
		return e
	}

	// output params
	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return err
		}
	}

	return nil
}`

// Generate writes the Go client implementations to w.
func Generate(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf

	out(w, "// Client is the API client.\n")
	out(w, "type Client struct {\n")
	out(w, "  // URL is the required API endpoint address.\n")
	out(w, "  URL string\n\n")
	out(w, "  // AuthToken is an optional authentication token.\n")
	out(w, "  AuthToken string\n\n")
	out(w, "  // HTTPClient is the client used for making requests, defaulting to http.DefaultClient.\n")
	out(w, "  HTTPClient *http.Client\n")
	out(w, "}\n\n")

	for _, m := range s.Methods {
		name := format.GoName(m.Name)
		out(w, "// %s %s\n", name, m.Description)
		out(w, "func (c *Client) %s(", name)

		// input arg
		if len(m.Inputs) > 0 {
			out(w, "in %sInput", name)
		}
		out(w, ") ")

		// output arg
		if len(m.Outputs) > 0 {
			out(w, "(*%sOutput, error) {\n", name)
			out(w, "  var out %sOutput\n", name)
		} else {
			out(w, "error {\n")
		}

		// return
		out(w, "  return ")
		if len(m.Outputs) > 0 {
			out(w, "&out, ")
		}
		out(w, "call(c.HTTPClient, c.AuthToken, c.URL, \"%s\", ", m.Name)
		if len(m.Inputs) > 0 {
			out(w, "in, ")
		} else {
			out(w, "nil, ")
		}
		if len(m.Outputs) > 0 {
			out(w, "&out)\n")
		} else {
			out(w, "nil)\n")
		}

		// close
		out(w, "}\n\n")
	}

	out(w, "\n%s\n", call)

	return nil
}
