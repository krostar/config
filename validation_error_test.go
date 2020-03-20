package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_String(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		var e = make(ValidationError)

		assert.Equal(t, "no validation errors", e.String())
	})

	t.Run("one error", func(t *testing.T) {
		var e = make(ValidationError)

		e["one"] = errors.New("errone")
		assert.Equal(t, "validation error: field one errone", e.String())
	})

	t.Run("two errors", func(t *testing.T) {
		var e = make(ValidationError)

		e["one"] = errors.New("errone")
		e["two"] = errors.New("errtwo")

		s := e.String()

		assert.True(t, strings.HasPrefix(s, "validation error: "))
		assert.Contains(t, s, "field one errone")
		assert.Contains(t, s, "field two errtwo")
	})
}
