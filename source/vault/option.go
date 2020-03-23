package sourcevault

import (
	"path/filepath"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Option defines the function signature to apply options.
type Option func(v *Vault)

// SecretMayNotExist makes calls to SetValueFromConfigTreePath
// function to return an error that implements IsTrivial when
// file does not exists.
func SecretMayNotExist() Option {
	return func(v *Vault) { v.strict = false }
}

// WithClient provides a vault client.
func WithClient(client *vaultapi.Client) Option {
	return func(v *Vault) { v.client = client }
}

// WithCustomTreePathMatcher provides a custom matcher
// between the configuration tree path and the vault key path.
func WithCustomTreePathMatcher(matcher TreePathMatcherFunc) Option {
	return func(v *Vault) { v.treePathMatcher = matcher }
}

// DefaultTreePathMatcher defines a default matcher using the provided
// prefix and the config path tree.
// For example, with a prefix 'my_app', and a config tree 'database.postgres.password'
// it will return the path 'my_app/database/postgres' and the key 'password'.
func DefaultTreePathMatcher(prefix string) Option {
	return func(v *Vault) {
		v.treePathMatcher = func(treePath string) (string, string) {
			treePath = strings.ReplaceAll(treePath, ".", "/")

			pathWithoutKey := filepath.Dir(treePath)
			if pathWithoutKey == "." || pathWithoutKey == "/" {
				pathWithoutKey = ""
			}

			keyWithoutPath := filepath.Base(treePath)
			if keyWithoutPath == "." || keyWithoutPath == "/" {
				keyWithoutPath = ""
			}

			if pathWithoutKey != "" {
				pathWithoutKey = "/" + pathWithoutKey
			}

			return prefix + pathWithoutKey, keyWithoutPath
		}
	}
}
