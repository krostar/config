package config_test

import (
	"fmt"
	"os"

	"github.com/krostar/config"
	sourceenv "github.com/krostar/config/source/env"
)

type Config struct {
	NoDefault string
	Some      OtherConfig
	SomeOther *OtherConfig
	SomeLast  *OtherConfig
}

type OtherConfig struct {
	UniversalAnswer int
}

func (c *OtherConfig) SetDefault() {
	c.UniversalAnswer = 42
}

func ExampleLoad() {
	var cfg Config

	os.Setenv("MYAPP_SOMEOTHER_UNIVERSALANSWER", "1010") // nolint: errcheck, gosec
	defer os.Unsetenv("MYAPP_SOMEOTHER_UNIVERSALANSWER") // nolint: errcheck, gosec

	if err := config.Load(&cfg, config.WithSources(
		sourceenv.New("myapp"),
	)); err != nil {
		panic(err)
	}

	fmt.Printf("cfg.NoDefault                 = %q\n", cfg.NoDefault)
	fmt.Printf("cfg.Some.UniversalAnswer      = %d\n", cfg.Some.UniversalAnswer)
	fmt.Printf("cfg.SomeOther.UniversalAnswer = %d\n", cfg.SomeOther.UniversalAnswer)
	fmt.Printf("cfg.SomeLast                  = %v\n", cfg.SomeLast)

	// Output:
	// cfg.NoDefault                 = ""
	// cfg.Some.UniversalAnswer      = 42
	// cfg.SomeOther.UniversalAnswer = 1010
	// cfg.SomeLast                  = <nil>
}

type HTTPConfig struct {
	Debug         bool
	ListenAddress string
}

func (c *HTTPConfig) SetDefault() {
	c.ListenAddress = "127.0.0.1:8080"
}

func (c HTTPConfig) Validate() error {
	if c.ListenAddress == "" {
		return fmt.Errorf("listening address can't be empty")
	}
	return nil
}

func ExampleValidate() {
	cfg := HTTPConfig{
		Debug:         false,
		ListenAddress: "",
	}

	err := config.Validate(&cfg)

	fmt.Println("failed:", err != nil)
	fmt.Println("reason:", err.Error())

	// Output:
	// failed: true
	// reason: validation error: listening address can't be empty
}

func ExampleSetDefault() {
	var cfg HTTPConfig

	if err := config.SetDefault(&cfg); err != nil {
		panic("unable to set defaults")
	}

	fmt.Println("debug:", cfg.Debug)
	fmt.Println("listen-address:", cfg.ListenAddress)

	// Output:
	// debug: false
	// listen-address: 127.0.0.1:8080
}
