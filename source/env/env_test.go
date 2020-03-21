package sourceenv

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/config/internal/trivialerr"
)

func newEnv(t *testing.T, prefix string) *Env {
	env, err := New(prefix)()
	require.NoError(t, err)
	return env.(*Env)
}

func TestEnv_SetValueFromConfigTreePath(t *testing.T) {
	const (
		prefix   = "prefix"
		prefixUp = "PREFIX"
	)

	var tests = map[string]struct {
		envKey                 string
		envValue               string
		key                    string
		expectedValue          interface{}
		expectedFailure        bool
		expectedTrivialFailure bool
	}{
		"test yolo": {
			key:                    "yolo",
			envKey:                 prefixUp + "_YOLO",
			envValue:               "42",
			expectedValue:          42,
			expectedFailure:        false,
			expectedTrivialFailure: false,
		}, "test yi li": {
			key:                    "yi_li",
			envKey:                 prefixUp + "_YI_LI",
			envValue:               "hello",
			expectedValue:          "hello",
			expectedFailure:        false,
			expectedTrivialFailure: false,
		}, "test not found": {
			key:                    "willnobefound",
			expectedValue:          "hello",
			expectedFailure:        true,
			expectedTrivialFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			env := newEnv(t, prefix)

			if test.envKey != "" {
				require.NoError(t, os.Setenv(test.envKey, test.envValue))
				defer func() {
					assert.NoError(t, os.Unsetenv(test.envKey))
				}()
			}

			v := reflect.New(reflect.TypeOf(test.expectedValue)).Elem()

			isset, err := env.SetValueFromConfigTreePath(&v, test.key)
			if test.expectedFailure {
				require.Error(t, err)
				if test.expectedTrivialFailure {
					require.True(t, trivialerr.IsTrivial(err))
				}
			} else {
				require.NoError(t, err)
				assert.True(t, isset)
				assert.Equal(t, test.expectedValue, v.Interface())
			}
		})
	}
}

func TestEnv_Name(t *testing.T) {
	require.Equal(t, "env", newEnv(t, "").Name())
}

func TestEnv_keyFormatter(t *testing.T) {
	tests := map[string]struct {
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

	env := newEnv(t, "prefix")

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			formattedKey := env.keyFormatter(test.keyName)
			assert.Equal(t, "PREFIX_"+test.expectedFormattedKey, formattedKey)
		})
	}
}
