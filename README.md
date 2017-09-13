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

Authentication is handled automatically via the API Client using the provided `username` and `password`.

If you would like to manually re-authenticate, you may do this by calling the `Authenticate` method on the `client`:

```go
client.Authenticate()
```

## Documentation

TODO: link to GoDoc

## Contributing

### Setup

This project uses [dep](https://github.com/golang/dep) for dependency management.

To install/update dep and all dependencies, run:

```bash
make setup
```

## Testing

```bash
make test
```

## TODO

- [ ] Finish unit tests and stub out JSON responses
- [ ] Support pagination
- [x] Auto-refresh token if expired before calling endpoints?
- [ ] Make sure all endpoints are covered
- [ ] Publish GoDoc
