package configue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithSources(t *testing.T) {
	var (
		c  Configue
		s1 Source
		s2 Source
	)

	assert.Empty(t, c.sources)
	WithSources(s1, s2)(&c)
	assert.Len(t, c.sources, 2)
}
