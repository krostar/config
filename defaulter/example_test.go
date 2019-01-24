package defaulter_test

import (
	"fmt"

	"github.com/krostar/configue/defaulter"
)

type HTTPConfig struct {
	Debug         bool
	ListenAddress string
}

func (c *HTTPConfig) SetDefault() {
	c.ListenAddress = "127.0.0.1:8080"
}

func Example() {
	var cfg HTTPConfig

	if err := defaulter.SetDefault(&cfg); err != nil {
		panic("unable to set defaults")
	}

	fmt.Println("debug:", cfg.Debug)
	fmt.Println("listen-address:", cfg.ListenAddress)

	// Output:
	// debug: false
	// listen-address: 127.0.0.1:8080
}
