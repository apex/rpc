
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
 * Call method with params via a POST request.
 */

async function call(url: string, method: string, authToken?: string, params?: any): Promise<string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  }
  
  if (authToken != null) {
    headers['Authorization'] = `Bearer ${authToken}`
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
   * Decoder is used as the reviver parameter when decoding responses.
   */

  private decoder(key: any, value: any) {
    return typeof value == 'string' && reISO8601.test(value)
      ? new Date(value)
      : value
  }

  /**
   * addItem: adds an item to the list.
   */

  async addItem(params: AddItemInput) {
    await call(this.url, 'add_item', this.authToken, params)
  }

  /**
   * getItems: returns all items in the list.
   */

  async getItems(): Promise<GetItemsOutput> {
    let res = await call(this.url, 'get_items', this.authToken)
    let out: GetItemsOutput = JSON.parse(res, this.decoder)
    return out
  }

  /**
   * removeItem: removes an item from the to-do list.
   */

  async removeItem(params: RemoveItemInput): Promise<RemoveItemOutput> {
    let res = await call(this.url, 'remove_item', this.authToken, params)
    let out: RemoveItemOutput = JSON.parse(res, this.decoder)
    return out
  }

}
