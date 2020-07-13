package tsclient

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/schema"
)

var require = `
// fetch for Node
const fetch = (typeof window == 'undefined' || window.fetch == null)
// @ts-ignore
  ? require('%s')
  : window.fetch
`

var call = `/**
 * ClientError is an API client error providing the HTTP status code and error type.
 */

class ClientError extends Error {
  status: number;
  type?: string;

  constructor(status: number, message?: string, type?: string) {
    super(message)
    this.status = status
    this.type = type
  }
}

/**
 * Call method with params via a POST request.
 */

async function call(url: string, method: string, authToken?: string, params?: any): Promise<string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  }
  
  if (authToken != null) {
    headers['Authorization'] = ` + "`Bearer ${authToken}`" + `
  }
  
  const res = await fetch(url + '/' + method, {
    method: 'POST',
    body: JSON.stringify(params),
    headers
  })

  // we have an error, try to parse a well-formed json
  // error response, otherwise default to status code
  if (res.status >= 300) {
    let err
    try {
      const { type, message } = await res.json()
      err = new ClientError(res.status, message, type)
    } catch {
      err = new ClientError(res.status, res.statusText)
    }
    throw err
  }

  return res.text()
}`

// Generate writes the TS client implementations to w.
func Generate(w io.Writer, s *schema.Schema, fetchLibrary string) error {
	out := fmt.Fprintf

	out(w, require, fetchLibrary)
	out(w, "\n%s\n", call)
	out(w, "\n\n")
	out(w, `const reISO8601 = /(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))/`)
	out(w, "\n\n")
	out(w, "/**\n")
	out(w, " * Client is the API client.\n")
	out(w, " */\n")
	out(w, "\n")
	out(w, "export class Client {\n")
	out(w, "\n")
	out(w, "  private url: string\n")
	out(w, "  private authToken?: string\n")
	out(w, "\n")
	out(w, "  /**\n")
	out(w, "   * Initialize.\n")
	out(w, "   */\n")
	out(w, "\n")
	out(w, "  constructor(params: { url: string, authToken?: string }) {\n")
	out(w, "    this.url = params.url\n")
	out(w, "    this.authToken = params.authToken\n")
	out(w, "  }\n")
	out(w, "\n")
	out(w, "  /**\n")
	out(w, "   * Decoder is used as the reviver parameter when decoding responses.\n")
	out(w, "   */\n")
	out(w, "\n")
	out(w, "  private decoder(key: any, value: any) {\n")
	out(w, "    return typeof value == 'string' && reISO8601.test(value)\n")
	out(w, "      ? new Date(value)\n")
	out(w, "      : value\n")
	out(w, "  }\n")
	out(w, "\n")

	// methods
	for _, m := range s.Methods {
		name := format.JsName(m.Name)
		out(w, "  /**\n")
		out(w, "   * %s: %s\n", name, m.Description)
		out(w, "   */\n\n")

		// input
		if len(m.Inputs) > 0 {
			out(w, "  async %s(params: %sInput)", name, format.GoName(m.Name))
		} else {
			out(w, "  async %s()", name)
		}

		// output
		if len(m.Outputs) > 0 {
			out(w, ": Promise<%sOutput> {\n", format.GoName(m.Name))
		} else {
			out(w, " {\n")
		}

		// return
		if len(m.Outputs) > 0 {
			out(w, "    let res = ")
			// call
			if len(m.Inputs) > 0 {
				out(w, "await call(this.url, '%s', this.authToken, params)\n", m.Name)
			} else {
				out(w, "await call(this.url, '%s', this.authToken)\n", m.Name)
			}
			out(w, "    let out: %sOutput = JSON.parse(res, this.decoder)\n", format.GoName(m.Name))
			out(w, "    return out\n")
		} else {
			// call
			if len(m.Inputs) > 0 {
				out(w, "    await call(this.url, '%s', this.authToken, params)\n", m.Name)
			} else {
				out(w, "    await call(this.url, '%s', this.authToken)\n", m.Name)
			}
		}

		out(w, "  }\n\n")
	}

	out(w, "}\n")

	return nil
}
