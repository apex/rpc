
v0.3.2 / 2020-09-01
===================

  * Revert "fix camel-casing of TS type field names" (forgot this will serialize wrong)

v0.3.1 / 2020-09-01
===================

  * fix camel-casing of TS type field names

v0.3.0 / 2020-09-01
===================

  * add dotnet client (#35)
  * add schema/schema.json to the binary. Closes #1
  * refactor ReadRequest() error response to mention a valid JSON object
  * update github workflow to use Go 1.14.x
  * remove Any type

v0.2.0 / 2020-08-06
===================

  * add exporting of TypeScript types
  * add start of Elm client (not usable yet)
  * add start of Ruby client. Closes #4
  * add start of PHP client. Closes #6
  * fix: remove inclusion of oneOf() util when no validation is present in gotypes
  * refactor test fixtures using tj/go-fixture

v0.1.2 / 2020-07-13
===================

  * revert camel-casing

v0.1.1 / 2020-07-13
===================

  * fix js field camel-casing

v0.1.0 / 2020-07-13
===================

  * add sorting of type fields
  * add support for custom fetch imports
  * fix TypeScript errors to support Deno
  * fix TypeScript output types, use interface not class
  * fix timestamp fields, now providing Date conversion
  * fix rpc-go-server to not use hard coded types package
  * fix enum support for optional fields
  * remove rpc-apex-docs, not useful to other people. Closes #16
