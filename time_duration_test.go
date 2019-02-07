package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomDuration_ToDuration(t *testing.T) {
	var (
		expectedDuration = time.Minute * 5
		customDuration   = customDuration(expectedDuration)
	)

	assert.Equal(t, expectedDuration, customDuration.ToDuration())
}

func TestCustomDuration_ToInt64(t *testing.T) {
	var (
		expectedInt    = int64(time.Minute * 5)
		customDuration = customDuration(expectedInt)
	)

	assert.Equal(t, expectedInt, customDuration.ToInt64())
}

func TestCustomDuration_MarshalJSON(t *testing.T) {
	var (
		customDuration = customDuration(time.Minute*5 + 45*time.Second)
		expectedRepr   = `"5m45s"`
	)

	repr, err := customDuration.MarshalJSON()

	require.NoError(t, err)
	assert.Equal(t, expectedRepr, string(repr))
}

func TestCustomDuration_UnmarshalJSON(t *testing.T) {
	var tests = map[string]struct {
		value            []byte
		expectedDuration time.Duration
		expectedFailure  bool
	}{
		"empty value": {
			value:           nil,
			expectedFailure: true,
		}, "float value": {
			value:            []byte("42000000000"),
			expectedDuration: 42 * time.Second,
		}, "valid string value": {
			value:            []byte("\"42s\""),
			expectedDuration: 42 * time.Second,
		}, "invalid string value": {
			value:           []byte("\"hello\""),
			expectedFailure: true,
		}, "invalid value type": {
			value:           []byte("[\"42s\"]"),
			expectedFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			var cd customDuration

			err := cd.UnmarshalJSON(test.value)
			if test.expectedFailure {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedDuration, cd.ToDuration())
		})
	}
}
