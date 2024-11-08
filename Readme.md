# Godi - Lightweight Modular Dependency Injection Framework for Go

`Godi` is a lightweight, modular dependency injection framework designed for building scalable Go applications. It simplifies the creation of highly maintainable Go applications through a modular architecture with clear boundaries and flexible configurations powered by Uber's Dig library under the hood.

## Key Features

- **Modular Architecture**: Organize your application into cohesive modules with well-defined interfaces.
- **Constructor-Based and Direct Injection**: Inject dependencies either via constructors or directly.
- **Type-Safe Dependency Management**: Ensure compile-time type safety in dependency injection.
- **Routing with Guards, Filters, Pipes, Interceptors**: Easily configure guards filters, interceptors for controllers or routes for flexible request handling.

## Getting Started

- Read the detailed documentation, please refer to the [Godi package documentation on pkg.go.dev](https://pkg.go.dev/github.com/huboh/godi#section-documentation).
- Add Godi to your project:

  ```bash
  go get github.com/huboh/godi
  ```

## Quick example

Here’s a minimal setup to start a godi app

```go
package main

import (
    "log"
    "github.com/huboh/godi"
)

const (
    port = "5000"
    host = "localhost"
)

func main() {
    app, err := godi.New(&main.module{})
    if err != nil {
        log.Fatal("Failed to create godi app:", err)
    }

    err := app.Listen(host, port);
    if err != nil {
        log.Fatal("Failed to start app server:", err)
    }
}
```

## Authors

- Knowledge - [Website](https://huboh.vercel.app)

## License

Godi is released under the [MIT License](https://github.com/huboh/godi/blob/main/LICENCE).
