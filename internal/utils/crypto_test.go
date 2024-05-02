package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandom(t *testing.T) {
	t.Run("generate random", func(t *testing.T) {
		const wantSize = 8
		got, err := GenerateRandom(wantSize)

		require.NoError(t, err)
		assert.Equal(t, wantSize, len(got))
	})
	t.Run("zero size", func(t *testing.T) {
		const wantSize = 0
		got, err := GenerateRandom(wantSize)

		require.NoError(t, err)
		assert.Equal(t, wantSize, len(got))
	})
	t.Run("size is less zero", func(t *testing.T) {
		const size = -1

		assert.Panics(t, func() {
			_, _ = GenerateRandom(size)
		})
	})
}

func TestGenerateRandomAES256Key(t *testing.T) {
	const wantKeySize = 32
	got, err := GenerateRandomAES256Key()

	require.NoError(t, err)
	assert.Equal(t, wantKeySize, len(got))
}

func TestBlockCipher_Seal(t *testing.T) {
	t.Run("seal and unseal", func(t *testing.T) {
		key, err := GenerateRandomAES256Key()
		require.NoError(t, err)
		sut := NewBlockCipher(key)
		plaintext := []byte("sensitive data")

		ciphertext, err := sut.Seal(plaintext)

		require.NoError(t, err)

		got, err := sut.Unseal(ciphertext)

		require.NoError(t, err)
		assert.Equal(t, plaintext, got)
	})
	t.Run("invalid key", func(t *testing.T) {
		const size = 8 // should be min 16 bytes
		key, err := GenerateRandom(size)
		require.NoError(t, err)
		sut := NewBlockCipher(key)
		plaintext := []byte("sensitive data")

		_, err = sut.Seal(plaintext)

		require.Error(t, err)
	})
	t.Run("seal empty data", func(t *testing.T) {
		key, err := GenerateRandomAES256Key()
		require.NoError(t, err)
		sut := NewBlockCipher(key)
		var plaintext []byte

		_, err = sut.Seal(plaintext)

		require.NoError(t, err)
	})
}

func TestBlockCipher_Unseal(t *testing.T) {
	t.Run("invalid key", func(t *testing.T) {
		const size = 8 // should be min 16 bytes
		key, err := GenerateRandom(size)
		require.NoError(t, err)
		sut := NewBlockCipher(key)
		ciphertext := []byte("ciphertext")

		_, err = sut.Unseal(ciphertext)

		require.Error(t, err)
	})
	t.Run("unseal invalid ciphertext", func(t *testing.T) {
		key, err := GenerateRandomAES256Key()
		require.NoError(t, err)
		sut := NewBlockCipher(key)
		ciphertext := []byte("ciphertext")

		_, err = sut.Unseal(ciphertext)

		require.Error(t, err)
	})
}
