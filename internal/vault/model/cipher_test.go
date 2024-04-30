package model

import (
	"crypto/aes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/utils"
)

const aes256KeySize = 2 * aes.BlockSize

func TestSeal(t *testing.T) {
	t.Run("seal and unseal data", func(t *testing.T) {
		const text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
		key, err := utils.GenerateRandom(aes256KeySize)
		require.NoError(t, err)
		want := &Secret{
			ID:   uuid.New(),
			Data: []byte(text),
		}
		sut := NewAESGCMCipher()

		sealed, err := sut.Seal(want, key)
		require.NoError(t, err)

		got, err := sut.Unseal(sealed, key)
		require.NoError(t, err)

		assert.Equal(t, want, got)
	})
	t.Run("seal with invalid key size", func(t *testing.T) {
		want := &Secret{
			ID:   uuid.New(),
			Data: []byte("data"),
		}
		sut := NewAESGCMCipher()
		const size = 8
		key, err := utils.GenerateRandom(size)
		require.NoError(t, err)

		_, err = sut.Seal(want, key)

		require.Error(t, err)
	})
	t.Run("seal empty data", func(t *testing.T) {
		key, err := utils.GenerateRandom(aes256KeySize)
		require.NoError(t, err)
		secret := &Secret{
			ID:   uuid.New(),
			Data: nil,
		}
		sut := NewAESGCMCipher()

		_, err = sut.Seal(secret, key)

		require.NoError(t, err)
	})
}

func TestUnseal(t *testing.T) {
	t.Run("unseal data with invalid key", func(t *testing.T) {
		secret := &Secret{
			ID:   uuid.New(),
			Data: []byte("data"),
		}
		sut := NewAESGCMCipher()
		const size = 8
		key, err := utils.GenerateRandom(size)
		require.NoError(t, err)

		_, err = sut.Unseal(secret, key)

		require.Error(t, err)
	})
	t.Run("unseal invalid data", func(t *testing.T) {
		secret := &Secret{
			ID:   uuid.New(),
			Data: []byte("invalid data"),
		}
		sut := NewAESGCMCipher()
		key, err := utils.GenerateRandom(aes256KeySize)
		require.NoError(t, err)

		_, err = sut.Unseal(secret, key)

		require.Error(t, err)
	})
}
