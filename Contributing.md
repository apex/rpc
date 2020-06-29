
## Contributing a language client

Contributing a client starts with creating a new [generator](./generators), a Go package which is tasked with creating the output for the new client, then a [command](./cmd) which wraps this package for use as a command-line tool.

Generated clients __MUST__:

- Set the Content-Type header field to `application/json`
- Support Authorization Bearer tokens
- Handle HTTP status errors reporting `>= 300` appropriately
- Support decoding of application/json error responses and exposing this information
- Suppor 204 "No Content" responses for methods which return no data

Generated clients __SHOULD__:

- Set the `User-Agent` header field with the name and version of the library
- Support compressed gzip requests
- Support compressed gzip responses