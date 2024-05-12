package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	vault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestCacheSecrets(t *testing.T) {
	t.Run("cache secrets", func(t *testing.T) {
		sut := New()
		want := []*vault.Secret{
			{ID: "1"},
			{ID: "2"},
		}

		sut.CacheSecrets(want)

		got := sut.ListSecrets()
		assert.ElementsMatch(t, want, got)

		sut.CacheSecrets(want)

		got = sut.ListSecrets()
		assert.ElementsMatch(t, want, got)
	})
	t.Run("cache empty secrets", func(t *testing.T) {
		sut := New()
		want := []*vault.Secret{}

		sut.CacheSecrets(want)

		got := sut.ListSecrets()
		assert.ElementsMatch(t, want, got)
	})
	t.Run("ignore already cached secret", func(t *testing.T) {
		want := []*vault.Secret{
			{ID: "1", Data: "data"},
			{ID: "2"},
		}
		secrets := []*vault.Secret{
			{ID: "1"},
			{ID: "2"},
		}
		sut := New()
		sut.CacheSecret(&vault.Secret{ID: "1", Data: "data"})

		sut.CacheSecrets(secrets)

		got := sut.ListSecrets()
		assert.ElementsMatch(t, want, got)
	})
}

func TestCacheSecret(t *testing.T) {
	t.Run("cache secret", func(t *testing.T) {
		sut := New()
		want := &vault.Secret{
			ID:   "1",
			Data: "secret",
		}

		sut.CacheSecret(want)

		got, ok := sut.GetSecret(want.ID)
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})
}
