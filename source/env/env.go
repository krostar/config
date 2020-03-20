// Package sourceenv sources configuration from env.
package sourceenv

import (
	"os"
	"strings"

	"github.com/krostar/config/internal/trivialerr"
)

// Env implements config.Source to fetch values from env
// based on the value's key.
type Env struct {
	prefix string
}

// New returns a new env source.
func New(prefix string) *Env {
	return &Env{prefix: prefix}
}

// Name implements config.Source interface.
func (e *Env) Name() string { return "env" }

func (e *Env) keyFormatter(key string) string {
	const sep = "_"
	return strings.
		NewReplacer(".", sep, "/", sep, "-", "", ",", "").
		Replace(strings.ToUpper(e.prefix + "_" + key))
}

// GetReprValueByKey gets the key's value from the system environment and return
// it. It return an error that implement IsTrivial when the key is not found.
// It never returns another kind of error.
func (e *Env) GetReprValueByKey(key string) ([]byte, error) {
	key = e.keyFormatter(key)

	env, exists := os.LookupEnv(key)
	if !exists {
		return nil, trivialerr.New("env does not contain key %s", key)
	}

	return []byte(env), nil
}
