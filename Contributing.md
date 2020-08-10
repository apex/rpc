
## Contributing a language client

Contributing a client starts with creating a new [generator](./generators), a Go package which is tasked with creating the output for the new client, then a [command](./cmd) which wraps this package for use as a command-line tool.

Generated clients __MUST__:

- Set the Content-Type header field to `application/json`
- Support Authorization Bearer tokens
- Handle HTTP status errors reporting `>= 300` appropriately
- Support decoding of application/json error responses and exposing this information
- Support 204 "No Content" responses for methods which return no data

Generated clients __SHOULD__:

- Set the `User-Agent` header field with the name and version of the library
- Support compressed gzip requests
- Support compressed gzip responses

## Testing

Clients use [tj/go-fixture](https://github.com/tj/go-fixture) for testing, use `go test -update` to generate or update any test fixtures.

## Schema

The JSON Schema used to validate Apex RPC's schema is located in the ./schema directory. Any changes made to this schema must be re-generated with `go generate`.