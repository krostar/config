package configue_test

import (
	"fmt"
	"os"

	"github.com/krostar/configue"
	sourceenv "github.com/krostar/configue/source/env"
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

func Example() {
	var cfg Config

	os.Setenv("MYAPP_SOMEOTHER_UNIVERSALANSWER", "1010") // nolint: errcheck, gosec

	if err := configue.Load(&cfg, configue.WithSources(
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
