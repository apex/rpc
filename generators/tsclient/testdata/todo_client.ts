
import fetch from 'node-fetch'

/**
* Call method with params via a POST request.
*/

async function call(url: string, authToken: string, method: string, params?: any): Promise<string> {
 const res = await fetch(url + '/' + method, {
	 method: 'POST',
	 body: JSON.stringify(params),
	 headers: {
		 'Content-Type': 'application/json',
		 'Authorization': `Bearer ${authToken}`
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
		 err = new Error(`${res.status} ${res.statusText}`)
	 }
	 throw err
 }

 return res.text()
}


const reDate = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2}(?:\.\d*))(?:Z|(\+|-)([\d|:]*))?$/

/**
 * Client is the API client.
 */

export class Client {

  private url: string
  private authToken: string

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
    return typeof value == 'string' && reDate.test(value)
      ? new Date(value)
      : value
  }

  /**
   * addItem: Add an item to the list.
   */

  async addItem(params: AddItemInput) {
    await call(this.url, this.authToken, 'add_item', params)
  }

  /**
   * getItems: Return all items in the list.
   */

  async getItems(): Promise<GetItemsOutput> {
    let res = await call(this.url, this.authToken, 'get_items')
    let out: GetItemsOutput = JSON.parse(res, this.decoder)
    return out
  }

  /**
   * removeItem: removes an item from the to-do list.
   */

  async removeItem(params: RemoveItemInput) {
    await call(this.url, this.authToken, 'remove_item', params)
  }

}
