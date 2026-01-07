# envalid

[![Release](https://img.shields.io/github/release/nyashmyash99/envalid.svg?style=flat-square)](https://github.com/nyashmyash99/envalid/releases)
[![Build Status](https://github.com/nyashmyash99/envalid/actions/workflows/envalid.yml/badge.svg?branch=master)](https://github.com/nyashmyash99/envalid/actions/workflows/envalid.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nyashmyash99/envalid)](https://goreportcard.com/report/github.com/nyashmyash99/envalid)
[![Go Reference](https://pkg.go.dev/badge/github.com/nyashmyash99/envalid)](https://pkg.go.dev/github.com/nyashmyash99/envalid)

**envalid** is a library for validating and accessing environment variables in Go, inspired by a [similar library in TypeScript](https://github.com/af/envalid).


## Quick start

```go
package config

import (
  env "github.com/nyashmyash99/envalid"
)

func ptr[T any](v T) *T { return &v }

type Config struct {
  Port      int 
  JwtSecret string
}

func Load() (*Config, error) {
  return env.Load[Config](env.Schema{
    "Port": env.Supplier(env.ParsePort, env.Variable[uint16]{
      Key:         "PORT",
      Description: "HTTP server port",
      Default:     ptr(uint16(8080)),
    }),
    "JwtSecret": env.Supplier(env.ParseStr, env.Variable[string]{
      Key: "JWT_SECRET",
    }),
  })
}
```


## Documentation

### Validator types

`ParseStr` - ensures that env exists.
> Note that an empty string is considered a valid value.

`ParseInt32` - ensures that env is an int32.

`ParseInt64` - ensures that env is an int64.

`ParseFloat32` - ensures that env is a float32.

`ParseFloat64` - ensures that env is a float64.

`ParseBool` - ensures that env is a bool-like.
> true: 1, true, t, yes, y, on

> false: 0, false, f, no, n, off

\* List is expanded by changing `BoolMap`.

`ParsePort` - ensures that env is a port (1-65535).

Not what you need? I welcome [contribution](https://github.com/NyashMyash99/envalid?tab=readme-ov-file#contribution).


### Custom validators

```go
// Create a parser.
func ParseAdmin(s string) (bool, error) {
  return s == "NyashMyash99", nil
}

// Ensure that the parser implements Parser.
var _ Parser[bool] = ParseAdmin

type Config struct {
  Admin bool
}

func Load() (*Config, error) {
  return env.Load[Config](env.Schema{
    // Use the parser with the appropriate Variable type.
    "Admin": env.Supplier(ParseAdmin, env.Variable[bool]{
      Key: "USER",
    }),
  })
}
```


## Contribution

I appreciate contributions!

Suggest improvements in [Issues](https://github.com/NyashMyash99/envalid/issues) or [Pull requests](https://github.com/NyashMyash99/envalid/pulls) using an idiomatic style of code/documentation and [Conventional Commits](https://www.conventionalcommits.org).


## License

This project is licensed under the [MIT License](https://github.com/NyashMyash99/envalid/blob/master/LICENSE) - see the LICENSE file for details.


## TODO

### Parsers

- Domain
- IPv4
- IPv6
- Host (domain + ipv4 + ipv6)
- URL
- DSN (URL)
- Email

### Options

- Variants - an array of available values for env.
