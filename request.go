package rpc

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ReadRequest parses application/json request bodies into value, or returns an error.
func ReadRequest(r *http.Request, value interface{}) error {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		// decode
		err := json.NewDecoder(r.Body).Decode(value)
		if err != nil {
			return BadRequest("Failed to parse malformed request body, must be a valid JSON object")
		}

		// validate
		if v, ok := value.(Validator); ok {
			err := v.Validate()
			if err != nil {
				return Invalid(err.Error())
			}
		}

		return nil
	default:
		return BadRequest("Unsupported request Content-Type, must be application/json")
	}
}
