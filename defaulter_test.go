package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type structComplex struct {
	Untouched    string
	Discard      stringDefaultable `default:"-"`
	Default      stringDefaultable
	PTRUntouched *stringDefaultable
	PTRTouched   *stringDefaultable
	Simple       structSimple
	unexported   stringDefaultable
}

type structSimple struct {
	Different stringDefaultable
	String    string
}

func (sc *structSimple) SetDefault() { sc.Different = "people" }

type stringDefaultable string

func (s *stringDefaultable) SetDefault() { *s = "world" }

func Test_SetDefault(t *testing.T) {
	sc := structComplex{
		PTRTouched: new(stringDefaultable),
		unexported: "",
	}
	defaultable := stringDefaultable("world")
	expectedSC := structComplex{
		PTRTouched: &defaultable,
		Default:    "world",
		Simple: structSimple{
			Different: "people",
		},
	}

	require.NoError(t, SetDefault(&sc))
	assert.Equal(t, expectedSC, sc)

	assert.Error(t, SetDefault(nil))
}

func Test_tryToSetDefault(t *testing.T) {
	var (
		zeroString                  = ""
		expectedZeroString          = zeroString
		nonZeroString               = "hello"
		expectedNonZeroString       = nonZeroString
		zeroCustomString            = stringDefaultable("")
		expectedZeroCustomString    = stringDefaultable("world")
		nonZerocustomString         = stringDefaultable("hello")
		expectedNonZeroCustomString = nonZerocustomString
	)

	tests := map[string]struct {
		value         reflect.Value
		expectedValue reflect.Value
	}{
		"invalid value": {
			value: reflect.ValueOf(nil),
		}, "zero string": {
			value:         reflect.ValueOf(&zeroString),
			expectedValue: reflect.ValueOf(&expectedZeroString),
		}, "non-zero string": {
			value:         reflect.ValueOf(&nonZeroString),
			expectedValue: reflect.ValueOf(&expectedNonZeroString),
		}, "zero defaultable string": {
			value:         reflect.ValueOf(&zeroCustomString),
			expectedValue: reflect.ValueOf(&expectedZeroCustomString),
		}, "non-zero defaultable string": {
			value:         reflect.ValueOf(&nonZerocustomString),
			expectedValue: reflect.ValueOf(&expectedNonZeroCustomString),
		},
	}

	for name, test := range tests {
		// rescope test as we're gonna run them in parallel
		var test = test
		t.Run(name, func(t *testing.T) {
			var value = reflect.Indirect(test.value)
			tryToSetDefault(&value)
			if test.expectedValue.IsValid() {
				assert.Equal(t, test.expectedValue.Interface(), test.value.Interface())
			}
		})
	}
}

func Test_isZeroValue(t *testing.T) {
	tests := map[string]struct {
		value              reflect.Value
		isExpectedToBeZero bool
	}{
		"invalid value": {
			value:              reflect.ValueOf(nil),
			isExpectedToBeZero: false,
		}, "empty string": {
			value:              reflect.ValueOf(""),
			isExpectedToBeZero: true,
		}, "non empty string": {
			value:              reflect.ValueOf("hello"),
			isExpectedToBeZero: false,
		}, "empty int": {
			value:              reflect.ValueOf(0),
			isExpectedToBeZero: true,
		}, "non empty int": {
			value:              reflect.ValueOf(42),
			isExpectedToBeZero: false,
		}, "empty struct": {
			value:              reflect.ValueOf(struct{ Name string }{}),
			isExpectedToBeZero: true,
		}, "non empty struct": {
			value:              reflect.ValueOf(struct{ Name string }{Name: "John Doe"}),
			isExpectedToBeZero: false,
		}, "empty array": {
			value:              reflect.ValueOf(([]string)(nil)),
			isExpectedToBeZero: true,
		}, "non empty array": {
			value:              reflect.ValueOf([]string{"hello", "world"}),
			isExpectedToBeZero: false,
		}, "empty map": {
			value:              reflect.ValueOf((map[string]string)(nil)),
			isExpectedToBeZero: true,
		}, "non empty map": {
			value:              reflect.ValueOf(map[string]string{"hello": "world"}),
			isExpectedToBeZero: false,
		},
	}

	for name, test := range tests {
		// rescope test as we're gonna run them in parallel
		var test = test
		t.Run(name, func(t *testing.T) {
			isZero := isZeroValue(&test.value)
			assert.Equal(t, test.isExpectedToBeZero, isZero)
		})
	}
}
