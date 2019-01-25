# configue

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/configue)
[![Licence](https://img.shields.io/github/license/krostar/configue.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/configue.svg?style=for-the-badge)
![Lastest version release date](https://img.shields.io/github/release-date/krostar/configue.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/configue/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/configue.svg?branch=master)
[![Code quality](https://img.shields.io/codacy/grade/14e0121b7ace47afa5022d5db6d0858c/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/configue/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/14e0121b7ace47afa5022d5db6d0858c.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/configue/dashboard)

A simple yet useful configuration package.

## Motivation

On any project I've made personnally or for a company, except if the project was really (really) small, I always needed at one point to be able to configure a component in the project (the http listening port, the database credential, a business configuration, ...). I've been using __viper__ for some times now, but I was not really happy about it for some reasons (usage of strings keys to get configuration, globally defined configuration which are a pita in big project to understand what's used where, and to use in test, ...). I also used __confita__ From my point of view a configuration package should:

- have a __priorization of "sources"__ (for example file < env < cli args)
- be __strongly typed__ (it should not use string keys, or return interface{})
- be __modulable__ (to add a new "source" to retrieve configuration from vauld or consul for example)
- handle __defaults values__ (without string keys, and as close as the configuration definition)
- have a __clear and easy to use API__
- be __light__
- encourage and follow the __best practices__

That's what I tried to do in this configuration package which is made of 3 components:

- the default setter (the `defaulter` package) which handles defaults
- the sources (anything that implements one of the two sources interfaces) responsible for the retrieval of the configuration
- the "loader" (the `configue` package) which is responsible to set the `default` if any and to call each `sources`

## Usage / examples

First thing first, lets load the configuration from a file:

```go
var sources = []configue.Source{
    sourcefile.NewSource(configFile),
}

configue.Load(&to, configue.WithSources(sources...)) // handle err

// cfg is now usable
}
```

Now if we wanted to add a env source that could override values defined before:

```go
var sources = []configue.Source{
    sourcefile.NewSource(configFile),
    sourceenv.NewSource("prefix"),
}
```

Now let's handle defaults:

```go
// let's define a structure that hold our http configuration, for example
type HTTPConfig struct {
    ListenAddress  string        `json:"listen-address"  yaml:"listen-address"`
    RequestTimeout time.Duration `json:"request-timeout" yaml:"request-timeout"`
    MACSecret      []byte        `json:"mac-secret"      yaml:"mac-secret"`
    TLS            *TLSConfig    `json:"tls"             yaml:"tls"`
}

// SetDefault sets sane default for http config.
func (c *HTTPConfig) SetDefault() {
    c.ListenAddress = ":8080"
    c.RequestTimeout = 3 * time.Second
}

func foo() {
    // ... say env source is already defined
    // and the program lauched with PREFIX_MACSECRET="secret"
    // and the program lauched with PREFIX_LISTENADDRESS=":8082"

    var cfg HTTPConfig
    configue.Load(&cfg, configue.WithSources(envsource)) // handle err

    // cfg.ListenAddress = ":8082"
    // cfg.RequestTimeout = "3s"
    // cfg.MACSecret = "secret"
    // cfg.TLS = nil
}
```

More doc and examples in the configue's [godoc](https://godoc.org/github.com/krostar/configue)

## License

This project is under the MIT licence, please see the LICENCE file.
