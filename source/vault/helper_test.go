package sourcevault

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	vaultapi "github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	hclog.DefaultOutput = ioutil.Discard
	os.Exit(m.Run())
}

type secret string

func (s secret) IsSecret() bool { return true }

func newTestVaultAPI(t *testing.T, secrets map[string]map[string]interface{}) (*vaultapi.Client, func()) {
	core, keyShares, token := vault.TestCoreUnsealed(t)
	_ = keyShares

	ln, addr := vaulthttp.TestServer(t, core)

	config := vaultapi.DefaultConfig()
	config.Address = addr

	client, err := vaultapi.NewClient(config)
	require.NoError(t, err)

	client.SetToken(token)

	for path, secret := range secrets {
		_, err := client.Logical().Write(path, secret)
		require.NoError(t, err)
	}

	return client, func() {
		require.NoError(t, ln.Close())
		require.NoError(t, core.Shutdown())
	}
}

func newVault(t *testing.T, opts ...Option) (*Vault, func()) {
	client, stop := newTestVaultAPI(t, map[string]map[string]interface{}{
		"secret/apps/my_super_app/database/postgres": {
			"username": "pg-user",
			"password": "pg-pwd",
		},
	})

	opts = append([]Option{
		WithClient(client),
		DefaultTreePathMatcher("secret/apps/my_super_app"),
	}, opts...)

	vault, err := New(opts...)
	require.NoError(t, err)

	return vault, stop
}
