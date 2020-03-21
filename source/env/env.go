// Package sourceenv sources configuration from env.
package sourceenv

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/krostar/config"
	"github.com/krostar/config/internal/trivialerr"
)

// Env implements config.Source to fetch values from env
// based on the value's key.
type Env struct {
	prefix string
}

// New returns a new env source.
func New(prefix string) config.SourceCreationFunc {
	return func() (config.Source, error) {
		return &Env{prefix: prefix}, nil
	}
}

// Name implements config.Source interface.
func (e *Env) Name() string { return "env" }

func (e *Env) keyFormatter(key string) string {
	const sep = "_"
	return strings.
		NewReplacer(".", sep, "/", sep, "-", "", ",", "").
		Replace(strings.ToUpper(e.prefix + "_" + key))
}

// SetValueFromConfigTreePath gets the key's value from the system environment and
// set it. It return an error that implement IsTrivial when the key is not found.
func (e *Env) SetValueFromConfigTreePath(v *reflect.Value, treePath string) (bool, error) {
	treePath = e.keyFormatter(treePath)

	env, exists := os.LookupEnv(treePath)
	if !exists {
		return false, trivialerr.New("env does not contain key %s", treePath)
	}

	newV, err := config.InitializeNewValueOfTypeWithString(v.Type(), env)
	if err != nil {
		return false, fmt.Errorf("unable to initialize new value from %q: %w", env, err)
	}

	return config.SetNewValue(v, newV)
}
