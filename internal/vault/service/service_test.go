package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/nestjam/goph-keeper/internal/vault/repository/inmemory"
)

func TestAddSecret(t *testing.T) {
	t.Run("add secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{Data: []byte("text")}
		userID := uuid.New()

		got, err := sut.AddSecret(ctx, secret, userID)

		require.NoError(t, err)
		_, err = secretRepo.GetSecret(ctx, got, userID)
		require.NoError(t, err)
	})
	t.Run("invalid data key", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		rootKey := randomMasterKey(t)
		cipher := model.NewMasterKeyCipher(rootKey)
		setInvalidDataKey(t, ctx, cipher, keyRepo)
		secretRepo := inmemory.NewSecretRepository()

		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{}
		userID := uuid.New()

		_, err := sut.AddSecret(ctx, secret, userID)

		require.Error(t, err)
	})
}

func TestUpdateSecret(t *testing.T) {
	t.Run("update secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		userID := uuid.New()
		secret := &model.Secret{}
		var err error
		secret.ID, err = sut.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		secret.Data = []byte("edited text")

		err = sut.UpdateSecret(ctx, secret, userID)

		require.NoError(t, err)
		got, err := sut.GetSecret(ctx, secret.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, secret, got)
	})
	t.Run("invalid data key", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		rootKey := randomMasterKey(t)
		cipher := model.NewMasterKeyCipher(rootKey)
		setInvalidDataKey(t, ctx, cipher, keyRepo)
		secretRepo := inmemory.NewSecretRepository()

		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{}
		userID := uuid.New()
		_, err := secretRepo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)

		err = sut.UpdateSecret(ctx, secret, userID)

		require.Error(t, err)
	})
}

func TestGetSecret(t *testing.T) {
	t.Run("get secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		userID := uuid.New()
		want := &model.Secret{Data: []byte("text")}
		var err error
		want.ID, err = sut.AddSecret(ctx, want, userID)
		require.NoError(t, err)

		got, err := sut.GetSecret(ctx, want.ID, userID)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("key not found", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		secret := &model.Secret{
			Data:  []byte("text"),
			KeyID: uuid.New(),
		}
		userID := uuid.New()
		var err error
		secret.ID, err = secretRepo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)

		_, err = sut.GetSecret(ctx, secret.ID, userID)

		require.Error(t, err)
	})
}

func TestDeleteSecret(t *testing.T) {
	t.Run("delete secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		userID := uuid.New()
		secret := &model.Secret{}
		var err error
		secret.ID, err = secretRepo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)

		err = sut.DeleteSecret(ctx, secret.ID, userID)

		require.NoError(t, err)
		_, err = secretRepo.GetSecret(ctx, secret.ID, userID)
		assert.ErrorIs(t, err, vault.ErrSecretNotFound)
	})
	t.Run("failed to delete secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := &secretRepositoryMock{
			DeleteSecretFunc: func(ctx context.Context, secretID, userID uuid.UUID) error {
				return errors.New("failed")
			},
		}
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		userID := uuid.New()
		secretID := uuid.New()

		err := sut.DeleteSecret(ctx, secretID, userID)

		require.Error(t, err)
	})
}

func TestListSecrets(t *testing.T) {
	t.Run("list user secrets", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		userID := uuid.New()
		s := &model.Secret{}
		s.ID, _ = secretRepo.AddSecret(ctx, s, userID)
		s2 := &model.Secret{}
		s2.ID, _ = secretRepo.AddSecret(ctx, s2, userID)
		user2ID := uuid.New()
		s3 := &model.Secret{}
		_, _ = secretRepo.AddSecret(ctx, s3, user2ID)
		want := uuid.UUIDs{s.ID, s2.ID}

		secrets, err := sut.ListSecrets(ctx, userID)

		require.NoError(t, err)
		got := make([]uuid.UUID, len(secrets))
		for i := 0; i < len(secrets); i++ {
			got[i] = secrets[i].ID
		}
		assert.ElementsMatch(t, want, got)
	})
	t.Run("failed to list secret", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := &secretRepositoryMock{
			ListSecretsFunc: func(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
				return nil, errors.New("failed")
			},
		}
		rootKey := randomMasterKey(t)
		sut := NewVaultService(secretRepo, keyRepo, rootKey)
		userID := uuid.New()

		_, err := sut.ListSecrets(ctx, userID)

		require.Error(t, err)
	})
}

func setInvalidDataKey(t *testing.T, ctx context.Context, cipher *model.MasterKeyCipher, r vault.DataKeyRepository) {
	t.Helper()

	const keySize = 8 // should be 32
	key, _ := utils.GenerateRandom(keySize)
	dataKey := &model.DataKey{Key: key}
	dataKey, err := cipher.Seal(dataKey)
	require.NoError(t, err)
	_, err = r.RotateKey(ctx, dataKey)
	require.NoError(t, err)
}
