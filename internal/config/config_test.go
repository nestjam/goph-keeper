package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	want := &Config{
		Server: ServerConfig{
			Address:  defaultServerAddress,
			CertFile: defaultCertFile,
			KeyFile:  defaultKeyFile,
		},
	}

	got, err := New()

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFromArgs(t *testing.T) {
	t.Run("shorthand arg name", func(t *testing.T) {
		args := []string{
			"app",
			"-a",
			"http://localhost:8081",
			"-k",
			"1234",
		}
		want := &Config{
			Server: ServerConfig{
				Address:  "http://localhost:8081",
				CertFile: defaultCertFile,
				KeyFile:  defaultKeyFile,
			},
			Vault: VaultConfig{
				MasterKey: "1234",
			},
		}

		got, err := New(FromArgs(args))

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("full arg name", func(t *testing.T) {
		args := []string{
			"app",
			"--server.address",
			"http://localhost:8080",
			"--vault.masterkey",
			"psw",
		}
		want := &Config{
			Server: ServerConfig{
				Address:  "http://localhost:8080",
				CertFile: defaultCertFile,
				KeyFile:  defaultKeyFile,
			},
			Vault: VaultConfig{
				MasterKey: "psw",
			},
		}

		got, err := New(FromArgs(args))

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("failed to parse flags", func(t *testing.T) {
		args := []string{
			"app",
			"--flag",
			"http://localhost:8080",
		}

		_, err := New(FromArgs(args))

		require.Error(t, err)
	})
}

func TestFromYaml(t *testing.T) {
	t.Run("read config from yaml", func(t *testing.T) {
		yml :=
			`postgres:
  dataSourceName: postgres://user:psw/db
`
		r := strings.NewReader(yml)
		want := &Config{
			Server: ServerConfig{
				Address:  defaultServerAddress,
				CertFile: defaultCertFile,
				KeyFile:  defaultKeyFile,
			},
			Postgres: PostgresConfig{
				DataSourceName: "postgres://user:psw/db",
			},
		}

		got, err := New(FromYaml(r))

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("invalid yaml file", func(t *testing.T) {
		yml :=
			`- postgres:
dataSourceName: postgres://user:psw/db
`
		r := strings.NewReader(yml)

		_, err := New(FromYaml(r))

		assert.Error(t, err)
	})
}
