package config

import (
	"github.com/pkg/errors"

	"github.com/krostar/config/defaulter"
	"github.com/krostar/config/trivialerr"
)

// Config stores the source configuration applied through options.
type Config struct {
	sources []Source
}

// New creates a new config instance configured through options.
func New(opts ...Option) *Config {
	var c Config

	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// Load creates a new instance (see New), and call the Load method of it (see Config.Load).
func Load(cfg interface{}, opts ...Option) error {
	return New(opts...).Load(cfg)
}

// Load tries to apply defaults to the provided interface (see the defaulter package) and
// call all sources to load the configuration.
func (c *Config) Load(cfg interface{}, opts ...Option) error {
	for _, opt := range opts {
		opt(c)
	}

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

func (c *Config) loadSource(source Source, cfg interface{}) error {
	var err error

	if s, ok := source.(SourceUnmarshal); ok {
		err = s.Unmarshal(cfg)
	} else if s, ok := source.(SourceGetReprValueByKey); ok {
		err = loadThroughReflection(s, cfg)
	} else {
		err = errors.Errorf("%s does not fulfill any load interface", source.Name())
	}

	if trivialerr.IsTrivial(errors.Cause(err)) {
		return nil
	}
	return err
}
