# config

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/config)
[![Licence](https://img.shields.io/github/license/krostar/config.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/config.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/config/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/config)
[![Code quality](https://img.shields.io/codacy/grade/4369c8e78a3e4fd995bac6b963c500b3/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/config/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/4369c8e78a3e4fd995bac6b963c500b3.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/config/dashboard)

A simple yet powerfull configuration package.

## Motivation

On any project I've made personnally or for a company, except if the project was really
(really) small, I always needed at one point to be able to configure a component in the
project (the http listening port, the database credentials, a business configuration, ...).
I've been using **viper** for some times now, but I was not really happy about it for some
reasons (usage of strings keys to get configuration, globally defined configuration which are
a pita in big project to understand what's used where, and to concurently use and modify in test, ...). I also used **confita** from which this project was inspired.

From my point of view a configuration package should:

-   have a **priorization of "sources"** (for example file &lt; env &lt; cli args)
-   be **strongly typed** (it should not use string keys, or return interface{})
-   be **modulable** (to add a new "source" to retrieve
        configuration from vauld or consul for example)
-   handle **defaults values** (without string keys, and as close as the configuration definition)
-   have a **clear and easy to use API**
-   be **light**
-   encourage and follow the **best practices**

That's what I tried to do in this configuration package which is made of 3 components:

-   the default setter (the `defaulter` package) which handles defaults
-   the sources (anything that implements one of the two sources interfaces) responsible for
        the retrieval of the configuration
-   the "loader" (the `config` package) which is responsible to set the `default` if any
        and to call each `sources`

## Usage / examples

```go
// let's define a structure that hold our http configuration, for example
type HTTPConfig struct {
    Debug          bool
    ListenAddress  string
    RequestTimeout time.Duration 
    MACSecret      []byte
}

// SetDefault sets sane default for http config.
func (c *HTTPConfig) SetDefault() {
    c.ListenAddress = ":8080"
    c.RequestTimeout = 3 * time.Second
}

// Validate checks whenever the config is properly set.
func (c *HTTPConfig) Validate() error {
    if c.RequestTimeout < time.Second {
        return errors.New("request timeout is too short (min 1s)")
    }
}

func main() {
    // export PREFIX_DEBUG="true"
    // export PREFIX_MACSECRET="secret"
    // echo "{ "listen-address": ":8082" }" > ./conf.json

    var cfg HTTPConfig

    if err := config.Load(&cfg, config.WithSources(config.Source{
        sourcefile.New("./conf.json"),
        sourceenv.New("prefix"),
    })); err != nil {
        panic(err)
    }

    if err := config.Validate(&cfg); err != nil {
        panic(err)
    }

    // cfg.Debug          = "true"
    // cfg.ListenAddress  = ":8082"
    // cfg.RequestTimeout = "3s"
    // cfg.MACSecret      = "secret"
}
```

More doc and examples in the config's [godoc](https://godoc.org/github.com/krostar/config)

## License

This project is under the MIT licence, please see the LICENCE file.
