package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/nestjam/goph-keeper/internal/vault/repository/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddSecret(t *testing.T) {
	t.Run("add secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey, _ := utils.GenerateRandomAES256Key()
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{Data: []byte("text")}
		userID := uuid.New()

		got, err := sut.AddSecret(ctx, secret, userID)

		require.NoError(t, err)
		stored, err := secretRepo.GetSecret(ctx, got.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, stored.ID, got.ID)
		assert.NotNil(t, stored.IV)
		dataKey, _ := keyRepo.GetKey(ctx)
		assert.Equal(t, dataKey.ID, got.KeyID)
	})
	t.Run("invalid data key", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		const keySize = 8 // should be 32
		key, _ := utils.GenerateRandom(keySize)
		dataKey := &model.DataKey{Key: key}
		_, err := keyRepo.RotateKey(ctx, dataKey)
		require.NoError(t, err)
		secretRepo := inmemory.NewSecretRepository()
		rootKey, _ := utils.GenerateRandomAES256Key()
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{}
		userID := uuid.New()

		_, err = sut.AddSecret(ctx, secret, userID)

		require.Error(t, err)
	})
}

func TestGetSecret(t *testing.T) {
	t.Run("get secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey, _ := utils.GenerateRandomAES256Key()
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{Data: []byte("text")}
		wantData := secret.Data
		userID := uuid.New()
		added, err := sut.AddSecret(ctx, secret, userID)
		wantID := added.ID
		require.NoError(t, err)

		got, err := sut.GetSecret(ctx, added.ID, userID)

		require.NoError(t, err)
		assert.Equal(t, wantID, got.ID)
		assert.Equal(t, wantData, got.Data)
		assert.Nil(t, got.IV)
		assert.Equal(t, uuid.Nil, got.KeyID)
	})
	t.Run("key not found", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey, _ := utils.GenerateRandomAES256Key()
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{
			Data:  []byte("text"),
			KeyID: uuid.New(),
		}
		userID := uuid.New()
		added, err := secretRepo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)

		_, err = sut.GetSecret(ctx, added.ID, userID)

		require.Error(t, err)
	})
}
