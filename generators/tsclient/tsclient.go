package tsclient

import (
	"fmt"
	"io"

	"github.com/apex/rpc/internal/format"
	"github.com/apex/rpc/schema"
)

var call = `import fetch from 'node-fetch'

/**
* Call method with params via a POST request.
*/

async function call(url: string, authToken: string, method: string, params?: any): Promise<string> {
 const res = await fetch(url + '/' + method, {
	 method: 'POST',
	 body: JSON.stringify(params),
	 headers: {
		 'Content-Type': 'application/json',
		 'Authorization': ` + "`Bearer ${authToken}`" + `
	 }
 })

 // we have an error, try to parse a well-formed json
 // error response, otherwise default to status code
 if (res.status >= 300) {
	 let err
	 try {
		 const { type, message } = await res.json()
		 err = new Error(message)
		 err.type = type
	 } catch {
		 err = new Error(` + "`${res.status} ${res.statusText}`" + `)
	 }
	 throw err
 }

 return res.text()
}`

// Generate writes the TS client implementations to w.
func Generate(w io.Writer, s *schema.Schema) error {
	out := fmt.Fprintf

	out(w, "\n%s\n", call)
	out(w, "\n\n")
	out(w, "/**\n")
	out(w, " * Client is the API client.\n")
	out(w, " */\n")
	out(w, "\n")
	out(w, "export class Client {\n")
	out(w, "\n")
	out(w, "  private url: string\n")
	out(w, "  private authToken: string\n")
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
				out(w, "await call(this.url, this.authToken, '%s', params)\n", m.Name)
			} else {
				out(w, "await call(this.url, this.authToken, '%s')\n", m.Name)
			}
			out(w, "    let out: %sOutput = JSON.parse(res)\n", format.GoName(m.Name))
			out(w, "    return out\n")
		} else {
			// call
			if len(m.Inputs) > 0 {
				out(w, "    await call(this.url, this.authToken, '%s', params)\n", m.Name)
			} else {
				out(w, "    await call(this.url, this.authToken, '%s')\n", m.Name)
			}
		}

		out(w, "  }\n\n")
	}

	out(w, "}\n")

	return nil
}
