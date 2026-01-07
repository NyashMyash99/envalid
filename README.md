<h1 align="center">envalid</h1>

<p align="center">
  <strong>envalid</strong> is a library for validating and accessing environment variables in Go, inspired by a <a href="https://github.com/af/envalid" target="_blank">similar library in TypeScript</a>.
</p>

<p align="center">
  <a href="https://github.com/nyashmyash99/envalid/releases">
    <img src="https://img.shields.io/github/release/nyashmyash99/envalid.svg?style=flat-square" />
  </a>

  <a href="https://github.com/nyashmyash99/envalid/actions/workflows/envalid.yml">
    <img src="https://github.com/nyashmyash99/envalid/actions/workflows/envalid.yml/badge.svg?branch=master" />
  </a>

  <a href="https://codecov.io/github/NyashMyash99/envalid">
    <img src="https://codecov.io/github/NyashMyash99/envalid/graph/badge.svg?token=OMKEDW9BW2" />
  </a>

  <a href="https://goreportcard.com/report/github.com/nyashmyash99/envalid">
    <img src="https://goreportcard.com/badge/github.com/nyashmyash99/envalid" />
  </a>

  <a href="https://pkg.go.dev/github.com/nyashmyash99/envalid">
    <img src="https://pkg.go.dev/badge/github.com/nyashmyash99/envalid" />
  </a>
</p>

<hr />


## ⚡️ Quick start

```go
package config

import (
  env "github.com/nyashmyash99/envalid"
)

type Config struct {
  Port        int
  DatabaseUrl string
}

func Load() (*Config, error) {
  return env.Load[Config](env.Schema{
    "Port": env.Supplier(env.Port, env.Variable[uint16]{
      Key:         "PORT",
      Description: "HTTP server port",
      Default:     env.WithDefault(uint16(8080)),
    }),
    // postgres://user:password@localhost:5432/database
    "DatabaseUrl": env.Supplier(env.Str, env.Variable[string]{
      Key: "DATABASE_URL",
    }),
  })
}
```


## 📖 Documentation

### Validator types

`Str` - ensures that env exists.
> Note that an empty string is considered a valid value.

`Int32` - ensures that env is an int32.

`Int64` - ensures that env is an int64.

`Float32` - ensures that env is a float32.

`Float64` - ensures that env is a float64.

`Bool` - ensures that env is a bool-like.
> true: 1, true, t, yes, y, on

> false: 0, false, f, no, n, off

\* List is expanded by changing `BoolMap`.

`Port` - ensures that env is a port (1-65535).

Not what you need? I welcome [contribution](https://github.com/NyashMyash99/envalid?tab=readme-ov-file#contributing).


### Custom validators

```go
// Create a parser.
func ParseAdmin(s string) (bool, error) {
  return s == "NyashMyash99", nil
}

// Ensure that the parser implements Parser.
var _ Parser[bool] = ParseAdmin

// Shorten the function name by converting it to a validator, 
// or simply continue using the parser.
var Admin = Validator[bool](ParseAdmin)

type Config struct {
  Admin bool
}

func Load() (*Config, error) {
  return env.Load[Config](env.Schema{
    // Use the validator with the appropriate Variable type.
    "Admin": env.Supplier(Admin, env.Variable[bool]{
      Key: "USER",
    }),
  })
}
```


## 🤝 Contributing

I appreciate contributions!

Check out [contributing guidelines](https://github.com/nyashmyash99/envalid/blob/master/CONTRIBUTING.md) to learn more.


## 🔧 TODO

### Validators

- Domain
- IPv4
- IPv6
- Host (domain + ipv4 + ipv6)
- URL
- DSN (URL)
- Email

### Options

- Variants - an array of available values for env.
