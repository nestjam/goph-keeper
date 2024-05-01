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
		sut, err := NewDataKey()
		sut.ID = uuid.New()
		require.NoError(t, err)
		want := &Secret{
			ID:   uuid.New(),
			Data: []byte(text),
		}

		sealed, err := sut.Seal(want)

		require.NoError(t, err)
		assert.Equal(t, sut.ID, sealed.KeyID)
		assert.NotNil(t, sealed.IV)

		unsealed, err := sut.Unseal(sealed)

		require.NoError(t, err)
		assert.Equal(t, uuid.Nil, unsealed.KeyID)
		assert.Nil(t, unsealed.IV)

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
		sut := DataKey{Key: key}

		_, err = sut.Seal(want)

		require.Error(t, err)
	})
	t.Run("seal empty data", func(t *testing.T) {
		sut, err := NewDataKey()
		require.NoError(t, err)
		secret := &Secret{
			ID:   uuid.New(),
			Data: nil,
		}

		_, err = sut.Seal(secret)

		require.NoError(t, err)
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
		sut := &DataKey{Key: key}

		_, err = sut.Unseal(secret)

		require.Error(t, err)
	})
	t.Run("unseal invalid data", func(t *testing.T) {
		secret := &Secret{
			ID:   uuid.New(),
			Data: []byte("invalid data"),
		}
		key, err := NewDataKey()
		require.NoError(t, err)

		_, err = key.Unseal(secret)

		require.Error(t, err)
	})
}
