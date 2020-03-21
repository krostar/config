package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_InitializeNewValueOfTypeWithJSON(t *testing.T) {
	var tests = map[string]struct {
		valueRepr       []byte
		valueType       reflect.Type
		expectedValue   interface{}
		expectedFailure bool
	}{
		"nil values": {
			expectedFailure: true,
		},
		"create string value": {
			valueRepr:       []byte(`"hello"`),
			valueType:       reflect.TypeOf(""),
			expectedValue:   "hello",
			expectedFailure: false,
		},
		"create int value": {
			valueRepr:       []byte("42"),
			valueType:       reflect.TypeOf(0),
			expectedValue:   42,
			expectedFailure: false,
		},
		"create map value": {
			valueRepr:       []byte(`{ "key": "value" }`),
			valueType:       reflect.TypeOf(map[string]string{}),
			expectedValue:   map[string]string{"key": "value"},
			expectedFailure: false,
		},
		"create slice value": {
			valueRepr:       []byte(`["a", "b", "c", "d"]`),
			valueType:       reflect.TypeOf([]string{}),
			expectedValue:   []string{"a", "b", "c", "d"},
			expectedFailure: false,
		},
		"create pointer on string": {
			valueRepr: []byte(`"hello"`),
			// reflect.TypeOf("")   ==> get the type of a string
			// reflect.New().Addr() ==> create a new string, get a pointor on this string
			valueType:       reflect.TypeOf(reflect.New(reflect.TypeOf("")).Interface()),
			expectedValue:   "hello",
			expectedFailure: false,
		},
		"create time duration in float repr": {
			valueRepr:       []byte("1000000"),
			valueType:       reflect.TypeOf(time.Duration(0)),
			expectedValue:   time.Millisecond,
			expectedFailure: false,
		},
		"create time duration in string repr": {
			valueRepr:       []byte(`"17m"`),
			valueType:       reflect.TypeOf(time.Duration(0)),
			expectedValue:   17 * time.Minute,
			expectedFailure: false,
		},
		"create time duration in bad repr": {
			valueRepr:       []byte("10l"),
			valueType:       reflect.TypeOf(time.Duration(0)),
			expectedFailure: true,
		},
		"create int value with wrong type": {
			valueRepr:       []byte("string don't fit in int"),
			valueType:       reflect.TypeOf(0),
			expectedFailure: true,
		},
		"create empty in a int value": {
			valueType:       reflect.TypeOf(0),
			expectedFailure: true,
		},
		"create empty in a string value": {
			valueRepr:       []byte(`""`),
			valueType:       reflect.TypeOf(""),
			expectedValue:   "",
			expectedFailure: false,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value, err := InitializeNewValueOfTypeWithJSON(test.valueType, test.valueRepr)
			if test.expectedFailure {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedValue, reflect.Indirect(*value).Interface())
			}
		})
	}
}

func Test_InitializeNewValueOfTypeWithString(t *testing.T) {
	var tests = map[string]struct {
		valueRepr       string
		valueType       reflect.Type
		expectedValue   interface{}
		expectedFailure bool
	}{
		"nil values": {
			expectedFailure: true,
		},
		"create string value": {
			valueRepr:       "hello",
			valueType:       reflect.TypeOf(""),
			expectedValue:   "hello",
			expectedFailure: false,
		},
		"create int value": {
			valueRepr:       "42",
			valueType:       reflect.TypeOf(0),
			expectedValue:   42,
			expectedFailure: false,
		},
		"create map value": {
			valueRepr:       `{ "key": "value" }`,
			valueType:       reflect.TypeOf(map[string]string{}),
			expectedValue:   map[string]string{"key": "value"},
			expectedFailure: false,
		},
		"create slice value": {
			valueRepr:       `["a", "b", "c", "d"]`,
			valueType:       reflect.TypeOf([]string{}),
			expectedValue:   []string{"a", "b", "c", "d"},
			expectedFailure: false,
		},
		"create pointer on string": {
			valueRepr: "hello",
			// reflect.TypeOf("")   ==> get the type of a string
			// reflect.New().Addr() ==> create a new string, get a pointor on this string
			valueType:       reflect.TypeOf(reflect.New(reflect.TypeOf("")).Interface()),
			expectedValue:   "hello",
			expectedFailure: false,
		},
		"create time duration in float repr": {
			valueRepr:       "1000000",
			valueType:       reflect.TypeOf(time.Duration(0)),
			expectedValue:   time.Millisecond,
			expectedFailure: false,
		},
		"create time duration in string repr": {
			valueRepr:       "17m",
			valueType:       reflect.TypeOf(time.Duration(0)),
			expectedValue:   17 * time.Minute,
			expectedFailure: false,
		},
		"create time duration in bad repr": {
			valueRepr:       "10l",
			valueType:       reflect.TypeOf(time.Duration(0)),
			expectedFailure: true,
		},
		"create int value with wrong type": {
			valueRepr:       "string don't fit in int",
			valueType:       reflect.TypeOf(0),
			expectedFailure: true,
		},
		"create empty in a int value": {
			valueType:       reflect.TypeOf(0),
			expectedFailure: true,
		},
		"create empty in a string value": {
			valueRepr:       "",
			valueType:       reflect.TypeOf(""),
			expectedValue:   "",
			expectedFailure: false,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value, err := InitializeNewValueOfTypeWithString(test.valueType, test.valueRepr)
			if test.expectedFailure {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedValue, reflect.Indirect(*value).Interface())
			}
		})
	}
}

func Test_setValuesForEachAttributes(t *testing.T) {
	type (
		icfgNested struct {
			withUnexportedField int
			WithCustomTag       string `cfg:"olleh"`
			UniversalStrAnswer  string
		}

		icfg struct {
			WithDiscardTag            string `cfg:"-"`
			WithDiscardTagWithDefault string `cfg:"-"`
			WithCustomTag             string `cfg:"hello"`
			UniversalAnswer           int
			withUnexportedField       int
			WithDefaultValue          string
			WithUntouchedDefaultValue string
			WithNestedStruct          icfgNested
			WithInterface             interface{}
			WithNilPointer            *string
			WithStructPointer         *icfgNested
			WithStructNilPointer      *icfgNested
			WithPointer               *int
		}
	)
	var (
		nested = &icfgNested{
			UniversalStrAnswer: "lala",
		}
		i      = 42
		str    = "I'm not nil anymore"
		source = stubSourceThatUseReflection{
			"hello":                                   "world",
			"universalanswer":                         "42",
			"withunexportedfield":                     "-30",
			"withdefaultvalue":                        "I've been replaced",
			"withnestedstruct.withunexportedfield":    "-30",
			"withnestedstruct.olleh":                  "dlrow",
			"withnestedstruct.universalstranswer":     "forty-two",
			"withinterface":                           "blih",
			"withnilpointer":                          "I'm not nil anymore",
			"withpointer":                             "42",
			"withstructpointer.universalstranswer":    "lala",
			"withstructnilpointer.universalstranswer": "lala",
		}
		cfg = icfg{
			WithDiscardTagWithDefault: "I should NOT be replaced",
			withUnexportedField:       30,
			WithDefaultValue:          "I should be replaced",
			WithUntouchedDefaultValue: "I should NOT be replaced",
			WithPointer:               new(int),
			WithStructPointer:         new(icfgNested),
		}
		expectedCfg = icfg{
			WithDiscardTag:            "",
			WithDiscardTagWithDefault: "I should NOT be replaced",
			WithCustomTag:             "world",
			UniversalAnswer:           42,
			withUnexportedField:       30,
			WithDefaultValue:          "I've been replaced",
			WithUntouchedDefaultValue: "I should NOT be replaced",
			WithNestedStruct: icfgNested{
				withUnexportedField: 0,
				WithCustomTag:       "dlrow",
				UniversalStrAnswer:  "forty-two",
			},
			WithInterface:        "blih",
			WithNilPointer:       &str,
			WithPointer:          &i,
			WithStructPointer:    nested,
			WithStructNilPointer: nested,
		}
	)

	err := setValuesForEachAttributes(source, &cfg)
	require.NoError(t, err)
	assert.Equal(t, expectedCfg, cfg)
}
