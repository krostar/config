package sourcevault

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/config"
	"github.com/krostar/config/internal/trivialerr"
)

type configurationPostgres struct {
	Username secret
	Password secret
}

type configurationDatabase struct {
	Backend  string
	Postgres configurationPostgres
}

type configuration struct {
	Database configurationDatabase
	Other    string
}

func TestVault(t *testing.T) {
	vault, shutdown := newVault(t)
	defer shutdown()

	var cfg configuration
	err := config.Load(&cfg, config.WithRawSources(vault))
	require.NoError(t, err)
}

func TestVault_New(t *testing.T) {
	t.Run("treePathMatcher must be set", func(t *testing.T) {
		client, shutdown := newTestVaultAPI(t, nil)
		defer shutdown()

		_, err := New(WithClient(client))
		require.Error(t, err)
	})

	t.Run("client must be set", func(t *testing.T) {
		_, err := New(DefaultTreePathMatcher(""))
		require.Error(t, err)
	})

	t.Run("nominal with option", func(t *testing.T) {
		client, shutdown := newTestVaultAPI(t, nil)
		defer shutdown()

		vault, err := New(
			DefaultTreePathMatcher(""),
			WithClient(client),
		)
		require.NoError(t, err)
		require.NotNil(t, vault)
	})
}

func TestVault_Name(t *testing.T) {
	vault, clean := newVault(t)
	defer clean()
	assert.Equal(t, "vault", vault.Name())
}

func TestVault_SetValueFromConfigTreePath(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		vault, clean := newVault(t)
		defer clean()

		notASecret := reflect.New(reflect.TypeOf("")).Elem()
		pgUsername := reflect.New(reflect.TypeOf(secret(""))).Elem()
		pgPassword := reflect.New(reflect.TypeOf(secret(""))).Elem()
		kafkaPassword := reflect.New(reflect.TypeOf(secret(""))).Elem()

		update, err := vault.SetValueFromConfigTreePath(&notASecret, "database.postgres.username")
		require.NoError(t, err)
		assert.False(t, update)

		update, err = vault.SetValueFromConfigTreePath(&pgUsername, "database.postgres.username")
		require.NoError(t, err)
		assert.Equal(t, "pg-user", string(pgUsername.Interface().(secret)))
		assert.True(t, update)

		update, err = vault.SetValueFromConfigTreePath(&pgPassword, "database.postgres.password")
		require.NoError(t, err)
		assert.Equal(t, "pg-pwd", string(pgPassword.Interface().(secret)))
		assert.True(t, update)

		update, err = vault.SetValueFromConfigTreePath(&kafkaPassword, "database.kafka.password")
		require.Error(t, err)
		assert.False(t, update)
	})

	t.Run("nominal with SecretMayNotExist", func(t *testing.T) {
		vault, clean := newVault(t, SecretMayNotExist())
		defer clean()

		kafkaPassword := reflect.New(reflect.TypeOf(secret(""))).Elem()
		update, err := vault.SetValueFromConfigTreePath(&kafkaPassword, "database.kafka.password")
		require.True(t, trivialerr.IsTrivial(err))
		assert.False(t, update)
	})
}
