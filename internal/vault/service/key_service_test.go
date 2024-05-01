package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/nestjam/goph-keeper/internal/vault/repository/inmemory"
)

func TestKeyService_Seal(t *testing.T) {
	config := NewKeyRotationConfig()

	t.Run("seal secret data", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		key := setKey(t, ctx, keyRepo)
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{Data: []byte("data")}

		got, err := sut.Seal(ctx, secret)

		require.NoError(t, err)
		assert.Equal(t, secret.ID, got.ID)
		assert.Equal(t, key.ID, got.KeyID)
	})
	t.Run("data key is not set", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{Data: []byte("data")}

		got, err := sut.Seal(ctx, secret)

		require.NoError(t, err)
		assert.Equal(t, secret.ID, got.ID)
		key, _ := keyRepo.GetKey(ctx)
		assert.Equal(t, key.ID, got.KeyID)
	})
	t.Run("key rotation failed", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := &keyRepositoryMock{
			GetKeyFunc: func(ctx context.Context) (*model.DataKey, error) {
				return nil, vault.ErrKeyNotFound
			},
			RotateKeyFunc: func(ctx context.Context, key *model.DataKey) (*model.DataKey, error) {
				return nil, errors.New("failed")
			},
		}
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{Data: []byte("data")}

		_, err := sut.Seal(ctx, secret)

		require.Error(t, err)
	})
	t.Run("failed to get key", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := &keyRepositoryMock{
			GetKeyFunc: func(ctx context.Context) (*model.DataKey, error) {
				return nil, errors.New("failed")
			},
		}
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{Data: []byte("data")}

		_, err := sut.Seal(ctx, secret)

		require.Error(t, err)
	})
	t.Run("rotate key after n bytes are encrypted", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		config.EncryptedDataSizeThreshold = 5
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{Data: []byte("12345")} // 5 bytes
		sealed, err := sut.Seal(ctx, secret)
		require.NoError(t, err)
		keyID := sealed.KeyID

		secret = &model.Secret{Data: []byte("")} // 0 bytes
		sealed, err = sut.Seal(ctx, secret)
		require.NoError(t, err)
		key2ID := sealed.KeyID

		assert.NotEqual(t, keyID, key2ID)
		key, _ := keyRepo.GetKey(ctx)
		assert.Equal(t, key.ID, key2ID)
	})
	t.Run("rotate key after n encryptions are done", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		config.EncryptionsCountThreshold = 1
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{}
		sealed, err := sut.Seal(ctx, secret)
		require.NoError(t, err)
		keyID := sealed.KeyID

		secret = &model.Secret{}
		sealed, err = sut.Seal(ctx, secret)
		require.NoError(t, err)
		key2ID := sealed.KeyID

		assert.NotEqual(t, keyID, key2ID)
		key, _ := keyRepo.GetKey(ctx)
		assert.Equal(t, key.ID, key2ID)
	})
	t.Run("failed to update key stats", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := &keyRepositoryMock{
			GetKeyFunc: func(ctx context.Context) (*model.DataKey, error) {
				return nil, vault.ErrKeyNotFound
			},
			RotateKeyFunc: func(ctx context.Context, key *model.DataKey) (*model.DataKey, error) {
				return key, nil
			},
			UpdateStatsFunc: func(ctx context.Context, id uuid.UUID, dataSize int64) error {
				return errors.New("failed")
			},
		}
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{Data: []byte("data")}

		_, err := sut.Seal(ctx, secret)

		require.Error(t, err)
	})
}

func TestKeyService_Unseal(t *testing.T) {
	config := NewKeyRotationConfig()

	t.Run("unseal secret data", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		_ = setKey(t, ctx, keyRepo)
		sut := NewKeyService(keyRepo, config)
		want := &model.Secret{Data: []byte("data")}
		secret, err := sut.Seal(ctx, want)
		require.NoError(t, err)

		got, err := sut.Unseal(ctx, secret)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("key not found by id", func(t *testing.T) {
		ctx := context.Background()
		keyRepo := inmemory.NewDataKeyRepository()
		sut := NewKeyService(keyRepo, config)
		secret := &model.Secret{KeyID: uuid.New()}

		_, err := sut.Unseal(ctx, secret)

		require.Error(t, err)
	})
}

func setKey(t *testing.T, ctx context.Context, keyRepo vault.DataKeyRepository) *model.DataKey {
	t.Helper()

	key, _ := model.NewDataKey()
	key, err := keyRepo.RotateKey(ctx, key)
	require.NoError(t, err)

	return key
}
