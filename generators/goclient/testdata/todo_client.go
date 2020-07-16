// Client is the API client.
type Client struct {
  // URL is the required API endpoint address.
  URL string

  // AuthToken is an optional authentication token.
  AuthToken string

  // HTTPClient is the client used for making requests, defaulting to http.DefaultClient.
  HTTPClient *http.Client
}

// AddItem adds an item to the list.
func (c *Client) AddItem(in AddItemInput) error {
  return call(c.HTTPClient, c.AuthToken, c.URL, "add_item", in, nil)
}

// GetItems returns all items in the list.
func (c *Client) GetItems() (*GetItemsOutput, error) {
  var out GetItemsOutput
  return &out, call(c.HTTPClient, c.AuthToken, c.URL, "get_items", nil, &out)
}

// RemoveItem removes an item from the to-do list.
func (c *Client) RemoveItem(in RemoveItemInput) (*RemoveItemOutput, error) {
  var out RemoveItemOutput
  return &out, call(c.HTTPClient, c.AuthToken, c.URL, "remove_item", in, &out)
}


// Error is an error returned by the client.
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
}
