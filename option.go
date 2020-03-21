package config

// Option defines the function signature to apply options.
type Option func(c *Config) error

// WithRawSources appends the given initializd source to the list of configuration sources.
func WithRawSources(s ...Source) Option {
	return func(c *Config) error {
		c.sources = append(c.sources, s...)
		return nil
	}
}

// WithSources appends the given source to the list of configuration sources.
func WithSources(sf ...SourceCreationFunc) Option {
	return func(c *Config) error {
		for _, f := range sf {
			s, err := f()
			if err != nil {
				return err
			}
			c.sources = append(c.sources, s)
		}
		return nil
	}
}
