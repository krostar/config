package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithSources(t *testing.T) {
	var (
		c  Config
		s1 Source
		s2 Source
	)

	WithSources(s1, s2)(&c)
	assert.Len(t, c.sources, 2)
	assert.Equal(t, []Source{s1, s2}, c.sources)
}

func TestWithSourcesPrepend(t *testing.T) {
	var (
		c  Config
		s1 = dumbSource{}
		s2 = dumbSource{}
	)

	WithSources(s1)(&c)
	WithSourcesPrepend(s2)(&c)
	assert.Len(t, c.sources, 2)
	assert.Equal(t, []Source{s2, s1}, c.sources)
}
