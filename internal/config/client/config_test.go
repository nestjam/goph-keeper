package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	want := &Config{
		ServerAddress: "",
	}

	got, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFromArgs(t *testing.T) {
	t.Run("shorthand arg name", func(t *testing.T) {
		args := []string{
			"app",
			"-s",
			"http://localhost:8081",
		}
		want := &Config{
			ServerAddress: "http://localhost:8081",
		}

		got, err := NewConfig(FromArgs(args))

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("full arg name", func(t *testing.T) {
		args := []string{
			"app",
			"--serveraddress",
			"http://localhost:8080",
		}
		want := &Config{
			ServerAddress: "http://localhost:8080",
		}

		got, err := NewConfig(FromArgs(args))

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("failed to parse flags", func(t *testing.T) {
		args := []string{
			"app",
			"--flag",
			"http://localhost:8080",
		}

		_, err := NewConfig(FromArgs(args))

		require.Error(t, err)
	})
}
