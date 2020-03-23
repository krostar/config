// Package sourcevault sources configuration from vault.
package sourcevault

import (
	"errors"
	"fmt"
	"reflect"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/krostar/config"
	"github.com/krostar/config/internal/trivialerr"
)

// Vault implements config.Source with vault to fetch values
// based on the config key path.
type Vault struct {
	client          *vaultapi.Client
	strict          bool
	treePathMatcher TreePathMatcherFunc
}

// TreePathMatcherFunc defines the prototype of the tree path matcher function.
type TreePathMatcherFunc func(treePath string) (path, key string)

type secreter interface {
	IsSecret() bool
}

// New returns a new vault source.
func New(opts ...Option) (*Vault, error) {
	v := &Vault{
		strict: true,
	}

	for _, opt := range opts {
		opt(v)
	}

	if v.treePathMatcher == nil {
		return nil, fmt.Errorf("no tree path matcher configured")
	}

	if v.client == nil {
		return nil, fmt.Errorf("no vault client configured")
	}

	return v, nil
}

// Name implements config.Source interface.
func (v *Vault) Name() string { return "vault" }

// SetValueFromConfigTreePath tries to fetch the key from vault.
// It return an error that implement IsTrivial when the key is not found.
func (v *Vault) SetValueFromConfigTreePath(o *reflect.Value, treePath string) (bool, error) {
	if value, ok := o.Interface().(secreter); !ok || !value.IsSecret() {
		return false, nil
	}

	path, key := v.treePathMatcher(treePath)

	data, err := v.client.Logical().Read(path)
	if err != nil {
		return false, fmt.Errorf("unable to read vault secret %q: %w", treePath, err)
	}

	if data == nil {
		return false, trivialerr.WrapIf(v.strict, errors.New("vault secret not found"))
	}

	secret, ok := data.Data[key].(string)
	if !ok {
		return false, fmt.Errorf("secret for %q took from vault is not a string: %w", treePath, err)
	}

	newO, err := config.InitializeNewValueOfTypeWithString(o.Type(), secret)
	if err != nil {
		return false, fmt.Errorf("unable to initialize new value for %q: %w", treePath, err)
	}

	return config.SetNewValue(o, newO)
}
