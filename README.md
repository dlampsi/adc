# adc

[![Tests](https://github.com/dlampsi/adc/actions/workflows/tests.yml/badge.svg)](https://github.com/dlampsi/adc/actions/workflows/tests.yml)
[![Linter](https://github.com/dlampsi/adc/actions/workflows/linter.yml/badge.svg)](https://github.com/dlampsi/adc/actions/workflows/linter.yml)
[![codecov](https://codecov.io/gh/dlampsi/adc/graph/badge.svg?token=6TORMA0YJN)](https://codecov.io/gh/dlampsi/adc)
[![Go Reference](https://pkg.go.dev/badge/github.com/dlampsi/adc.svg)](https://pkg.go.dev/github.com/dlampsi/adc)

Active Directory client library that allows you to perform basic operations with users and groups: creation, deletion, search, changes to members and composition in groups.

The library is a wrapper around  [go-ldap/ldap](https://github.com/go-ldap/ldap) module that provides a more convient client for Active Directory.

## Usage

Import module in your go app:

```go
import "github.com/dlampsi/adc"
```

### Getting started

```go
cfg := &adc.Config{
    URL: "ldaps://my.ad.site:636",
    Bind: &adc.BindAccount{
        DN:       "CN=admin,DC=company,DC=com",
        Password: "***",
    },
    SearchBase: "OU=default,DC=company,DC=com",
}

cl := adc.New(cfg)

if err := cl.Connect(); err != nil {
    // Handle error
}

// Do stuff ...
```

See [examples](examples) directory for extended usage examples.

## Contributing

1. Fork this repositpry
2. Create new PR from `main` branch
2. Create PR from your fork
3. Make sure tests and coverage tests pass
4. Request review

## License

[MIT License](LICENSE).
