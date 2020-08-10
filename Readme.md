# RPC

Simple RPC style APIs with generated clients & servers.

## About

All RPC methods are invoked with the POST method, and the RPC method name is placed in the URL path. Input is passed as a JSON object in the body, following a JSON response for the output as shown here:

```sh
$ curl -d '{ "project_id": "ping_production" }' https://api.example.com/get_alerts
{
  "alerts": [...]
}
```

All inputs are objects, all outputs are objects, this improves future-proofing as additional fields can be added without breaking existing clients. This is similar to the approach AWS takes with their APIs.

## Commands

There are several commands provided for generating clients, servers, and documentation. Each of these commands accept a `-schema` flag defaulting to `schema.json`, see the `-h` help output for additional usage details.

### Clients

- `rpc-dotnet-client` generates .NET clients
- `rpc-ruby-client` generates Ruby clients
- `rpc-php-client` generates PHP clients
- `rpc-elm-client` generates Elm clients
- `rpc-go-client` generates Go clients
- `rpc-go-types` generates Go type definitions
- `rpc-ts-client` generates TypeScript clients

### Servers

- `rpc-go-server` generates Go servers

### Documentation

- `rpc-md-docs` generates markdown documentation

## Schemas

Currently the schemas are loosely a superset of [JSON Schema](https://json-schema.org/), however, this is a work in progress. See the [example schema](./examples/todo/schema.json).

## FAQ

<details>
  <summary>Why did you create this project?</summary>
  There are many great options when it comes to building APIs, but to me the most important aspect is simplicity, for myself and for the end user. Simple JSON in, and JSON out is appropriate for 99% of my API work, there's no need for the additional performance provided by alternative encoding schemes, and rarely a need for more complex features such as bi-directional streaming provided by gRPC.
</details>

<details>
  <summary>Should I use this in production?</summary>
  Only if you're confident that it supports everything you need, or you're comfortable with forking. I created this project for my work at Apex Software, it may not suit your needs.
</details>

<details>
  <summary>Why JSON schemas?</summary>
  I think concise schemas using a DSL are great, until they're a limiting factor. Personally I have no problem with JSON, and it's easy to expand upon when you introduce a new feature, such as inline examples for documentation.
</details>

<details>
  <summary>Why doesn't it follow the JSON-RPC spec?</summary>
  I would argue this spec is outdated, there is little reason to support batching at the request level, as HTTP/2 handles this for you.
</details>

<details>
  <summary>What does the client output look like?</summary>
  See the <a href="https://github.com/apex/logs/blob/master/client.go">Apex Logs</a> Go client for an example, client code is designed to be concise and idiomatic.
</details>

---

[![GoDoc](https://godoc.org/github.com/apex/rpc?status.svg)](https://godoc.org/github.com/apex/rpc)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

Sponsored by my [GitHub sponsors](https://github.com/sponsors/tj):

[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/0" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/0)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/1" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/1)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/2" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/2)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/3" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/3)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/4" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/4)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/5" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/5)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/6" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/6)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/7" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/7)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/8" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/8)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/9" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/9)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/10" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/10)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/11" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/11)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/12" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/12)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/13" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/13)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/14" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/14)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/15" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/15)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/16" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/16)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/17" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/17)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/18" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/18)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/19" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/19)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/20" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/20)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/21" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/21)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/22" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/22)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/23" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/23)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/24" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/24)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/25" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/25)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/26" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/26)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/27" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/27)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/28" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/28)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/29" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/29)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/30" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/30)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/31" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/31)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/32" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/32)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/33" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/33)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/34" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/34)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/35" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/35)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/36" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/36)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/37" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/37)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/38" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/38)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/39" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/39)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/40" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/40)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/41" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/41)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/42" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/42)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/43" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/43)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/44" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/44)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/45" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/45)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/46" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/46)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/47" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/47)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/48" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/48)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/49" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/49)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/50" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/50)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/51" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/51)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/52" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/52)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/53" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/53)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/54" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/54)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/55" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/55)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/56" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/56)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/57" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/57)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/58" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/58)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/59" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/59)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/60" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/60)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/61" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/61)
[<img src="https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/avatar/62" width="35">](https://sponsors-api-u2fftug6kq-uc.a.run.app/sponsor/profile/62)
