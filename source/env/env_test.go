package sourceenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/configue/trivialerr"
)

func TestEnv_GetReprValueByKey(t *testing.T) {
	const prefix = "prefix"
	const prefixUp = "PREFIX"

	var tests = map[string]struct {
		envKey                 string
		envValue               string
		key                    string
		expectedValue          []byte
		expectedFailure        bool
		expectedTrivialFailure bool
	}{
		"test yolo": {
			key:                    "yolo",
			envKey:                 prefixUp + "_YOLO",
			envValue:               "42",
			expectedValue:          []byte("42"),
			expectedFailure:        false,
			expectedTrivialFailure: false,
		}, "test yi li": {
			key:                    "yi_li",
			envKey:                 prefixUp + "_YI_LI",
			envValue:               "hello",
			expectedValue:          []byte("hello"),
			expectedFailure:        false,
			expectedTrivialFailure: false,
		}, "test not found": {
			key:                    "willnobefound",
			expectedValue:          []byte("hello"),
			expectedFailure:        true,
			expectedTrivialFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var env = New(prefix)

			if test.envKey != "" {
				require.NoError(t, os.Setenv(test.envKey, test.envValue))
				defer func() {
					assert.NoError(t, os.Unsetenv(test.envKey))
				}()
			}

			value, err := env.GetReprValueByKey(test.key)
			if test.expectedFailure {
				require.Error(t, err)
				if test.expectedTrivialFailure {
					require.True(t, trivialerr.IsTrivial(err))
				}
			} else {
				assert.Equal(t, test.expectedValue, value)
			}
		})
	}
}

func TestEnv_Name(t *testing.T) {
	require.Equal(t, "env", New("").Name())
}

func TestEnv_keyFormatter(t *testing.T) {
	var (
		env   = New("prefix")
		tests = map[string]struct {
			keyName              string
			expectedFormattedKey string
		}{
			"lowercase": {
				keyName:              "key",
				expectedFormattedKey: "KEY",
			}, "caps": {
				keyName:              "keyName",
				expectedFormattedKey: "KEYNAME",
			}, "dash": {
				keyName:              "key-name",
				expectedFormattedKey: "KEYNAME",
			}, "underscore": {
				keyName:              "key_name",
				expectedFormattedKey: "KEY_NAME",
			}, "dot": {
				keyName:              "key.name",
				expectedFormattedKey: "KEY_NAME",
			}, "slash": {
				keyName:              "key/name",
				expectedFormattedKey: "KEY_NAME",
			},
		}
	)

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			formattedKey := env.keyFormatter(test.keyName)
			assert.Equal(t, "PREFIX_"+test.expectedFormattedKey, formattedKey)
		})
	}
}
