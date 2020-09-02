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
 * Call method with body via a POST request.
 */

async function call(url: string, method: string, authToken?: string, body?: string): Promise<string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  }
  
  if (authToken != null) {
    headers['Authorization'] = ` + "`Bearer ${authToken}`" + `
  }
  
  const res = await fetch(url + '/' + method, {
    method: 'POST',
    body: body,
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
	out(w, "   * Decode the response to an object.\n")
	out(w, "   */\n")
	out(w, "\n")
	out(w, "  private decodeResponse(res: string): any {\n")
	out(w, "    const obj = JSON.parse(res)\n")
	out(w, "\n")
	out(w, "    const isObject = (val: any) =>\n")
	out(w, "      val && typeof val === \"object\" && val.constructor === Object\n")
	out(w, "    const isDate = (val: any) =>\n")
	out(w, "      typeof val == \"string\" && reISO8601.test(val)\n")
	out(w, "    const isArray = (val: any) => Array.isArray(val)\n")
	out(w, "\n")
	out(w, "    const decode = (val: any): any => {\n")
	out(w, "      let ret: any\n")
	out(w, "\n")
	out(w, "      if (isObject(val)) {\n")
	out(w, "        ret = {}\n")
	out(w, "        for (const prop in val) {\n")
	out(w, "          if (!Object.prototype.hasOwnProperty.call(val, prop)) {\n")
	out(w, "            continue\n")
	out(w, "          }\n")
	out(w, "          ret[this.toCamelCase(prop)] = decode(val[prop])\n")
	out(w, "        }\n")
	out(w, "      } else if (isArray(val)) {\n")
	out(w, "        ret = []\n")
	out(w, "        val.forEach((item: any) => {\n")
	out(w, "          ret.push(decode(item))\n")
	out(w, "        })\n")
	out(w, "      } else if (isDate(val)) {\n")
	out(w, "        ret = new Date(val)\n")
	out(w, "      } else {\n")
	out(w, "        ret = val\n")
	out(w, "      }\n")
	out(w, "\n")
	out(w, "      return ret\n")
	out(w, "    }\n")
	out(w, "\n")
	out(w, "    return decode(obj)\n")
	out(w, "  }\n")
	out(w, "\n")
	out(w, "  /**\n")
	out(w, "   * Convert a field name from snake case to camel case.\n")
	out(w, "   */\n")
	out(w, "\n")
	out(w, "  private toCamelCase(str: string): string {\n")
	out(w, "    const capitalize = (str: string) =>\n")
	out(w, "      str.charAt(0).toUpperCase() + str.slice(1)\n")
	out(w, "\n")
	out(w, "    const tok = str.split(\"_\")\n")
	out(w, "    let ret = tok[0]\n")
	out(w, "    tok.slice(1).forEach((t) => (ret += capitalize(t)))\n")
	out(w, "\n")
	out(w, "    return ret\n")
	out(w, "  }\n")
	out(w, "\n")
	out(w, "  /**\n")
	out(w, "   * Encode the request object.\n")
	out(w, "   */\n")
	out(w, "\n")
	out(w, "  private encodeRequest(obj: any): string {\n")
	out(w, "    const isObject = (val: any) =>\n")
	out(w, "      val && typeof val === \"object\" && val.constructor === Object\n")
	out(w, "    const isArray = (val: any) => Array.isArray(val)\n")
	out(w, "\n")
	out(w, "    const encode = (val: any): any => {\n")
	out(w, "      let ret: any\n")
	out(w, "\n")
	out(w, "      if (isObject(val)) {\n")
	out(w, "        ret = {}\n")
	out(w, "        for (const prop in val) {\n")
	out(w, "          if (!Object.prototype.hasOwnProperty.call(val, prop)) {\n")
	out(w, "            continue\n")
	out(w, "          }\n")
	out(w, "          ret[this.toSnakeCase(prop)] = encode(val[prop])\n")
	out(w, "        }\n")
	out(w, "      } else if (isArray(val)) {\n")
	out(w, "        ret = []\n")
	out(w, "        val.forEach((item: any) => {\n")
	out(w, "          ret.push(encode(item))\n")
	out(w, "        })\n")
	out(w, "      } else {\n")
	out(w, "        ret = val\n")
	out(w, "      }\n")
	out(w, "\n")
	out(w, "      return ret\n")
	out(w, "    }\n")
	out(w, "\n")
	out(w, "    return JSON.stringify(encode(obj))\n")
	out(w, "  }\n")
	out(w, "\n")
	out(w, "  /**\n")
	out(w, "   * Convert a field name from camel case to snake case.\n")
	out(w, "   */\n")
	out(w, "\n")
	out(w, "  private toSnakeCase(str: string): string {\n")
	out(w, "    let ret = \"\"\n")
	out(w, "    const isUpper = (c: string) => !(c.toLowerCase() == c)\n")
	out(w, "\n")
	out(w, "    for (let c of str) {\n")
	out(w, "      if (isUpper(c)) {\n")
	out(w, "        ret += \"_\" + c.toLowerCase()\n")
	out(w, "      } else {\n")
	out(w, "        ret += c\n")
	out(w, "      }\n")
	out(w, "    }\n")
	out(w, "\n")
	out(w, "    return ret\n")
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
				out(w, "await call(this.url, '%s', this.authToken, this.encodeRequest(params))\n", m.Name)
			} else {
				out(w, "await call(this.url, '%s', this.authToken)\n", m.Name)
			}
			out(w, "    let out: %sOutput = this.decodeResponse(res)\n", format.GoName(m.Name))
			out(w, "    return out\n")
		} else {
			// call
			if len(m.Inputs) > 0 {
				out(w, "    await call(this.url, '%s', this.authToken, this.encodeRequest(params))\n", m.Name)
			} else {
				out(w, "    await call(this.url, '%s', this.authToken)\n", m.Name)
			}
		}

		out(w, "  }\n\n")
	}

	out(w, "}\n")

	return nil
}
