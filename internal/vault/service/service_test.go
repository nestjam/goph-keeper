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
		repo := inmemory.NewSecretRepository()
		key, _ := utils.GenerateRandomAES256Key()
		sut := NewVaultService(repo, key)
		ctx := context.Background()
		secret := &model.Secret{Data: []byte("text")}
		userID := uuid.New()

		added, err := sut.AddSecret(ctx, secret, userID)

		require.NoError(t, err)
		stored, err := repo.GetSecret(ctx, added.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, added.ID, stored.ID)
		assert.NotNil(t, stored.IV)
	})
	t.Run("invalid aes key", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		const keySize = 8
		key, _ := utils.GenerateRandom(keySize)
		sut := NewVaultService(repo, key)
		ctx := context.Background()
		secret := &model.Secret{Data: []byte("text")}
		userID := uuid.New()

		_, err := sut.AddSecret(ctx, secret, userID)

		require.Error(t, err)
	})
}

func TestGetSecret(t *testing.T) {
	t.Run("get secret", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		key, _ := utils.GenerateRandomAES256Key()
		sut := NewVaultService(repo, key)
		ctx := context.Background()
		secret := &model.Secret{Data: []byte("text")}
		userID := uuid.New()
		added, err := sut.AddSecret(ctx, secret, userID)
		require.NoError(t, err)

		got, err := sut.GetSecret(ctx, added.ID, userID)

		require.NoError(t, err)
		assert.Equal(t, got.ID, added.ID)
		assert.Equal(t, secret.Data, got.Data)
		assert.Nil(t, got.IV)
	})
	t.Run("invalid sealed data", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		key, _ := utils.GenerateRandomAES256Key()
		sut := NewVaultService(repo, key)
		ctx := context.Background()
		secret := &model.Secret{Data: []byte("text")}
		userID := uuid.New()
		added, err := repo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)

		_, err = sut.GetSecret(ctx, added.ID, userID)

		require.Error(t, err)
	})
}
