package configue

import (
	"github.com/pkg/errors"

	"github.com/krostar/configue/defaulter"
	"github.com/krostar/configue/trivialerr"
)

// Configue stores the configuration applied through options.
type Configue struct {
	sources []Source
}

// New creates a new Configue instance configured through options.
func New(opts ...Option) *Configue {
	var c Configue

	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// Load creates a new instance (see New), and call the Load method of it (see Configue.Load).
func Load(cfg interface{}, opts ...Option) error {
	return New(opts...).Load(cfg)
}

// Load tries to apply defaults to the provided interface (see the defaulter package) and
// call all sources to load the configuration.
func (c *Configue) Load(cfg interface{}) error {
	if err := defaulter.SetDefault(cfg); err != nil {
		return errors.Wrap(err, "unable to set defaults")
	}

	for _, source := range c.sources {
		if err := c.loadSource(source, cfg); err != nil {
			return errors.Wrap(err, "failed to load configuration")
		}
	}
	return nil
}

func (c *Configue) loadSource(source Source, cfg interface{}) error {
	var err error

	if s, ok := source.(sourceThatUnmarshal); ok {
		err = s.Unmarshal(cfg)
	} else if s, ok := source.(sourceThatUseReflection); ok {
		err = reflectThroughConfig(s, cfg)
	} else {
		err = trivialerr.New("%s does not fulfill any load interface", source.Name())
	}

	if trivialerr.IsTrivial(errors.Cause(err)) {
		return nil
	}
	return err
}
