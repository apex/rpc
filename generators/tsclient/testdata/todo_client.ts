
// fetch for Node
const fetch = (typeof window == 'undefined' || window.fetch == null)
// @ts-ignore
  ? require('node-fetch')
  : window.fetch

/**
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
    headers['Authorization'] = `Bearer ${authToken}`
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
}


const reISO8601 = /(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d([+-][0-2]\d:[0-5]\d|Z))/

/**
 * Client is the API client.
 */

export class Client {

  private url: string
  private authToken?: string

  /**
   * Initialize.
   */

  constructor(params: { url: string, authToken?: string }) {
    this.url = params.url
    this.authToken = params.authToken
  }

  /**
   * Decode the response to an object.
   */

  private decodeResponse(res: string): any {
    const obj = JSON.parse(res)

    const isObject = (val: any) =>
      val && typeof val === "object" && val.constructor === Object
    const isDate = (val: any) =>
      typeof val == "string" && reISO8601.test(val)
    const isArray = (val: any) => Array.isArray(val)

    const decode = (val: any): any => {
      let ret: any

      if (isObject(val)) {
        ret = {}
        for (const prop in val) {
          if (!Object.prototype.hasOwnProperty.call(val, prop)) {
            continue
          }
          ret[this.toCamelCase(prop)] = decode(val[prop])
        }
      } else if (isArray(val)) {
        ret = []
        val.forEach((item: any) => {
          ret.push(decode(item))
        })
      } else if (isDate(val)) {
        ret = new Date(val)
      } else {
        ret = val
      }

      return ret
    }

    return decode(obj)
  }

  /**
   * Convert a field name from snake case to camel case.
   */

  private toCamelCase(str: string): string {
    const capitalize = (str: string) =>
      str.charAt(0).toUpperCase() + str.slice(1)

    const tok = str.split("_")
    let ret = tok[0]
    tok.slice(1).forEach((t) => (ret += capitalize(t)))

    return ret
  }

  /**
   * Encode the request object.
   */

  private encodeRequest(obj: any): string {
    const isObject = (val: any) =>
      val && typeof val === "object" && val.constructor === Object
    const isArray = (val: any) => Array.isArray(val)

    const encode = (val: any): any => {
      let ret: any

      if (isObject(val)) {
        ret = {}
        for (const prop in val) {
          if (!Object.prototype.hasOwnProperty.call(val, prop)) {
            continue
          }
          ret[this.toSnakeCase(prop)] = encode(val[prop])
        }
      } else if (isArray(val)) {
        ret = []
        val.forEach((item: any) => {
          ret.push(encode(item))
        })
      } else {
        ret = val
      }

      return ret
    }

    return JSON.stringify(encode(obj))
  }

  /**
   * Convert a field name from camel case to snake case.
   */

  private toSnakeCase(str: string): string {
    let ret = ""
    const isUpper = (c: string) => !(c.toLowerCase() == c)

    for (let c of str) {
      if (isUpper(c)) {
        ret += "_" + c.toLowerCase()
      } else {
        ret += c
      }
    }

    return ret
  }

  /**
   * addItem: adds an item to the list.
   */

  async addItem(params: AddItemInput) {
    await call(this.url, 'add_item', this.authToken, this.encodeRequest(params))
  }

  /**
   * getItems: returns all items in the list.
   */

  async getItems(): Promise<GetItemsOutput> {
    let res = await call(this.url, 'get_items', this.authToken)
    let out: GetItemsOutput = this.decodeResponse(res)
    return out
  }

  /**
   * removeItem: removes an item from the to-do list.
   */

  async removeItem(params: RemoveItemInput): Promise<RemoveItemOutput> {
    let res = await call(this.url, 'remove_item', this.authToken, this.encodeRequest(params))
    let out: RemoveItemOutput = this.decodeResponse(res)
    return out
  }

}
