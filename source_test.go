package config

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krostar/config/internal/trivialerr"
)

type stubSourceThatUseReflection map[string]string

func (s stubSourceThatUseReflection) Name() string { return "stub reflect" }

func (s stubSourceThatUseReflection) SetValueFromConfigTreePath(v *reflect.Value, name string) (bool, error) {
	var (
		str   string
		found bool
	)

	for key, resp := range s {
		if name == key {
			str = resp
			found = true
			break
		}
	}

	if !found {
		return false, nil
	}

	newV, err := InitializeNewValueOfTypeWithString(v.Type(), str)
	if err != nil {
		return false, fmt.Errorf("unable to initialize new value from %q: %w", str, err)
	}

	return SetNewValue(v, newV)
}

type stubSourceThatUnmarshal int

func (s stubSourceThatUnmarshal) Name() string { return "stub unmarshal" }

func (s stubSourceThatUnmarshal) Unmarshal(interface{}) error {
	var err error

	switch {
	case s < 0:
		err = trivialerr.New("not found")
	case s > 0:
		err = errors.New("not found")
	}

	return err
}

type dumbSource struct{}

func (dumbSource) Name() string { return "dumb" }

func Test_appendConfigTreePath(t *testing.T) {
	var tests = map[string]struct {
		parentName   string
		childName    string
		expectedName string
	}{
		"normal case": {
			parentName:   "parent",
			childName:    "child",
			expectedName: "parent.child",
		},
		"no parent": {
			parentName:   "",
			childName:    "child",
			expectedName: "child",
		},
		"uppercase": {
			parentName:   "ToTo",
			childName:    "tItI",
			expectedName: "toto.titi",
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			name := appendConfigTreePath(test.parentName, test.childName)
			assert.Equal(t, test.expectedName, name)
		})
	}
}
