package config

// Option defines the function signature to apply options.
type Option func(c *Config)

// WithSources appends the given source to the list of configuration sources.
func WithSources(s ...Source) Option {
	return func(c *Config) {
		c.sources = append(c.sources, s...)
	}
}

// WithSourcesPrepend prepends the given source to the list of configuration sources.
func WithSourcesPrepend(s ...Source) Option {
	return func(c *Config) {
		c.sources = append(s, c.sources...)
	}
}
