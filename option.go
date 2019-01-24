package configue

// Option defines the function signature to apply options.
type Option func(c *Configue)

// WithSources adds the given source to the list of configuration sources.
func WithSources(s ...Source) Option {
	return func(c *Configue) {
		c.sources = append(c.sources, s...)
	}
}
