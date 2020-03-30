package sourcevault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SecretMayNotExist(t *testing.T) {
	var v Vault

	v.strict = true
	SecretMayNotExist()(&v)
	assert.False(t, v.strict)
}

func Test_WithClient(t *testing.T) {
	client, shutdown := newTestVaultAPI(t, nil)
	defer shutdown()

	var v Vault

	WithClient(client)(&v)
	assert.NotNil(t, v.client)
}

func Test_WithCustomTreePathMatcher(t *testing.T) {
	var v Vault

	WithCustomTreePathMatcher(func(string) (string, string) { return "", "" })(&v)
	assert.NotNil(t, v.treePathMatcher)
}

func Test_DefaultTreePathMatcher(t *testing.T) {
	var v Vault

	DefaultTreePathMatcher("foo/bar")(&v)

	path, key := v.treePathMatcher("path1.path2.key")
	assert.Equal(t, "foo/bar/path1/path2", path)
	assert.Equal(t, "key", key)

	path, key = v.treePathMatcher("key")
	assert.Equal(t, "foo/bar", path)
	assert.Equal(t, "key", key)

	path, key = v.treePathMatcher("")
	assert.Equal(t, "foo/bar", path)
	assert.Equal(t, "", key)

	path, key = v.treePathMatcher(".")
	assert.Equal(t, "foo/bar", path)
	assert.Equal(t, "", key)

	path, key = v.treePathMatcher("..")
	assert.Equal(t, "foo/bar", path)
	assert.Equal(t, "", key)
}
