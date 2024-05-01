package vault

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DataKeyRepositoryContract struct {
	NewDataKeyRepository func() (DataKeyRepository, func())
}

func (c DataKeyRepositoryContract) Test(t *testing.T) {
	t.Run("add key", func(t *testing.T) {
		sut, tearDown := c.NewDataKeyRepository()
		t.Cleanup(tearDown)
		ctx := context.Background()
		key, err := model.NewDataKey()
		require.NoError(t, err)

		key, err = sut.RotateKey(ctx, key)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, key.ID)

		got, err := sut.GetKey(ctx)
		require.NoError(t, err)
		assert.Equal(t, key, got)

		got, err = sut.GetByID(ctx, key.ID)
		require.NoError(t, err)
		assert.Equal(t, key, got)
	})
	t.Run("rotate key", func(t *testing.T) {
		sut, tearDown := c.NewDataKeyRepository()
		t.Cleanup(tearDown)
		ctx := context.Background()
		key, err := model.NewDataKey()
		require.NoError(t, err)
		_, err = sut.RotateKey(ctx, key)
		require.NoError(t, err)
		key2, err := model.NewDataKey()
		require.NoError(t, err)

		key2, err = sut.RotateKey(ctx, key2)

		require.NoError(t, err)
		got, err := sut.GetKey(ctx)
		require.NoError(t, err)
		assert.Equal(t, key2, got)
	})
	t.Run("key not found by id", func(t *testing.T) {
		sut, tearDown := c.NewDataKeyRepository()
		t.Cleanup(tearDown)
		ctx := context.Background()
		id := uuid.New()

		_, err := sut.GetByID(ctx, id)

		require.ErrorIs(t, err, ErrKeyNotFound)
	})
	t.Run("active key no found", func(t *testing.T) {
		sut, tearDown := c.NewDataKeyRepository()
		t.Cleanup(tearDown)
		ctx := context.Background()

		_, err := sut.GetKey(ctx)

		require.ErrorIs(t, err, ErrKeyNotFound)
	})
	t.Run("update key stats", func(t *testing.T) {
		t.Run("update data size encrypted by key", func(t *testing.T) {
			sut, tearDown := c.NewDataKeyRepository()
			t.Cleanup(tearDown)
			ctx := context.Background()
			key, err := model.NewDataKey()
			require.NoError(t, err)
			key, err = sut.RotateKey(ctx, key)
			require.NoError(t, err)
			const want int64 = 100

			err = sut.UpdateStats(ctx, key.ID, want)

			require.NoError(t, err)
			key, err = sut.GetByID(ctx, key.ID)
			require.NoError(t, err)
			got := key.EncryptedSize
			assert.Equal(t, want, got)
		})
		t.Run("key not found", func(t *testing.T) {
			sut, tearDown := c.NewDataKeyRepository()
			t.Cleanup(tearDown)
			ctx := context.Background()
			key, err := model.NewDataKey()
			require.NoError(t, err)
			const want int64 = 100

			err = sut.UpdateStats(ctx, key.ID, want)

			require.ErrorIs(t, err, ErrKeyNotFound)
		})
	})
}
