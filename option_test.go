package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WithRawSources(t *testing.T) {
	var (
		c  Config
		s1 Source
		s2 Source
	)

	err := WithRawSources(s1, s2)(&c)
	require.NoError(t, err)

	assert.Len(t, c.sources, 2)
	assert.Equal(t, []Source{s1, s2}, c.sources)
}

func Test_WithSources(t *testing.T) {
	s1 := func() (Source, error) { return nil, nil }
	s2 := func() (Source, error) { return nil, errors.New("boum") }

	var c Config

	err := WithSources(s1, s1)(&c)
	require.NoError(t, err)
	assert.Len(t, c.sources, 2)

	err = WithSources(s1, s2)(&c)
	require.Error(t, err)
}
