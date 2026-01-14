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

type Author struct {
  Name  string
  Email string
}

type Config struct {
  Env         string
  HttpPort    uint16
  DatabaseUrl *env.PURL
  Authors     []Author
}

func Load() (*Config, error) {
  devDefDatabaseUrl, _ := env.ParseURL("postgres://local:password@localhost:5432/postgres?sslmode=disable")

  return env.Load[Config](env.Schema{
    "Env": env.Supplier(env.Str, env.Variable[string]{
      Key: "GO_ENV",
      Variants: []string{
        "development",
        "test",
        "staging",
        "production",
      },
      // Default: env.WithDefault("development"),
    }),
    "HttpPort": env.Supplier(env.Port, env.Variable[uint16]{
      Key:         "HTTP_PORT",
      Description: "HTTP server port",
      Default:     env.WithDefault(uint16(8080)),
    }),
    "DatabaseUrl": env.Supplier(env.URL, env.Variable[*env.PURL]{
      Key:        "DATABASE_URL",
      DevDefault: &devDefDatabaseUrl,
    }),
    "Authors": env.Supplier(env.JSON, env.Variable[[]Author]{
      Key: "AUTHORS",
    }),
  })
}
```


## 📖 Documentation

### Validator types

`Str` - ensures that the env var exists.
> Note that an empty string is considered a valid value.

`Int32`/`Int64`/`Float32`/`Float64` - ensures that the env var is a number.

`Bool` - ensures that the env var is a bool-like.
> true: 1, true, t, yes, y, on

> false: 0, false, f, no, n, off

\* List can be expanded by modifying `BoolMap`.

`Port` - ensures that the env var is a TCP/UDP port (1-65535).

`IPv4`/`IPv6` - ensures that the env var is a IP address.

`Domain` - ensures that the env var is a domain (including IDN).
> Minimal format: name.zone
 
> Maximal format: subdomain.name.zone

`Host` - ensures that the env var is a host (IPv4, IPv6, domain, including localhost).

`URL` - ensures that the env var is a url.
> Minimal format: protocol://host

> Maximal format: protocol://username:password@host:port/path?params#anchor

`JSON` - ensures that the env var is a JSON in the specified format.

`Email` - ensures that the env var is an email address in the format `user@domain`.

Not what you need? I welcome [contribution](https://github.com/NyashMyash99/envalid?tab=readme-ov-file#contributing).


### Validator options

`Description` - a string describing the env var.

`Variants` - an array of available values for the env var.
> Note that it case-sensitive.

`Default` - a fallback value that is returned if the env var has not been specified. Specifying a default value effectively makes the env var optional.
> Note that Default values takes precedence over validator and Variants, i.e. it may not match their conditions.

`DevDefault` - a fallback value that is returned if the env var has not been specified, GO_ENV var is specified and is not `production`.


### Custom validators

```go
package config

import (
	env "github.com/nyashmyash99/envalid"
)

// Create a validator.
var Admin = env.TrimmedParser(func(s string) (bool, error) {
  return s == "NyashMyash99", nil
})

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
