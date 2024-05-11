package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigFile(t *testing.T) {
	t.Run("default file", func(t *testing.T) {
		args := []string{}
		want := &ConfigFile{
			Name: defaultConfigFilePath,
		}

		got, err := NewConfigFile(args)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("from shorthand arg", func(t *testing.T) {
		args := []string{
			"app.exe",
			"-c",
			"./config.yml",
		}
		want := &ConfigFile{
			Name: "./config.yml",
		}

		got, err := NewConfigFile(args)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("from full arg", func(t *testing.T) {
		args := []string{
			"app.exe",
			"--config",
			"./my/config.yml",
		}
		want := &ConfigFile{
			Name: "./my/config.yml",
		}

		got, err := NewConfigFile(args)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
