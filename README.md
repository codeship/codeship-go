# Codeship API (v2) Client for Go

This is the start of an API client for the Codeship API written in Go.

[![Codeship Status for codeship/codeship-go](https://app.codeship.com/projects/c38f3280-792b-0135-21bb-4e0cf8ff365b/status?branch=master)](https://app.codeship.com/projects/244943)

## Usage

This library is intended to make integrating with Codeship fairly simple.

To start, you need to import the package:

```go
package main

import (
    codeship "github.com/codeship/codeship-go"
)
```

This library exposes the package `codeship`.

Getting a new API Client from it is done by calling `codeship.New()`:

```go
client := codeship.New("username", "password", "orgName")
```

## Authentication

Before performing any API requests you must first get an Authentication Token.

You do this by calling the `Authenticate` method on the `client`:

```go
client.Authenticate()
```

## Documentation

TODO: link to GoDoc

## Testing

```bash
make test
```

## TODO

- [ ] Finish unit tests and stub out JSON responses
- [ ] Support pagination
- [ ] (Optionally) Auto-refresh token if expired before calling endpoints?
- [ ] Make sure all endpoints are covered
- [ ] Publish GoDoc
