package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataKey_Copy(t *testing.T) {
	want, err := NewDataKey()
	want.ID = uuid.New()
	want.EncryptedDataSize = 1024
	want.EncryptionsCount = 4
	require.NoError(t, err)

	got := want.Copy()

	assert.Equal(t, want, got)
}

func TestDataKey_Seal(t *testing.T) {
	t.Run("seal and unseal data", func(t *testing.T) {
		const text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
		key, err := NewDataKey()
		key.ID = uuid.New()
		sut := NewDataKeyCipher(key)
		require.NoError(t, err)
		want := &Secret{
			ID:   uuid.New(),
			Data: []byte(text),
		}

		sealed, err := sut.Seal(want)

		require.NoError(t, err)
		assert.Equal(t, sut.dataKey.ID, sealed.KeyID)

		unsealed, err := sut.Unseal(sealed)

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil, unsealed.KeyID)

		assert.Equal(t, want, unsealed)
	})
	t.Run("seal with invalid key size", func(t *testing.T) {
		want := &Secret{
			ID:   uuid.New(),
			Data: []byte("data"),
		}
		const size = 8
		key, err := utils.GenerateRandom(size)
		require.NoError(t, err)
		dataKey := &DataKey{Key: key}
		sut := NewDataKeyCipher(dataKey)

		_, err = sut.Seal(want)

		require.Error(t, err)
	})
}

func TestDataKey_Unseal(t *testing.T) {
	t.Run("unseal data with invalid key", func(t *testing.T) {
		secret := &Secret{
			ID:   uuid.New(),
			Data: []byte("data"),
		}
		const size = 8
		key, err := utils.GenerateRandom(size)
		require.NoError(t, err)
		dataKey := &DataKey{Key: key}
		sut := NewDataKeyCipher(dataKey)

		_, err = sut.Unseal(secret)

		require.Error(t, err)
	})
}
