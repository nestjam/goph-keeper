package vault

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type SecretRepositoryContract struct {
	NewSecretRepository func() (SecretRepository, func())
}

func (c SecretRepositoryContract) Test(t *testing.T) {
	t.Run("add secret", func(t *testing.T) {
		sut, tearDown := c.NewSecretRepository()
		t.Cleanup(tearDown)
		secret := &model.Secret{}
		userID := uuid.New()
		ctx := context.Background()

		got, err := sut.AddSecret(ctx, secret, userID)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, got.ID)
	})
	t.Run("user has no secrets", func(t *testing.T) {
		sut, tearDown := c.NewSecretRepository()
		t.Cleanup(tearDown)
		userID := uuid.New()
		ctx := context.Background()

		got, err := sut.ListSecrets(ctx, userID)

		require.NoError(t, err)
		assert.Empty(t, got)
	})
	t.Run("list user secrets without sensitive data", func(t *testing.T) {
		sut, tearDown := c.NewSecretRepository()
		t.Cleanup(tearDown)
		userID := uuid.New()
		ctx := context.Background()
		secret := &model.Secret{Data: "1"}
		s1, err := sut.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		want := []*model.Secret{
			{ID: s1.ID},
		}

		got, err := sut.ListSecrets(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("list secrets of selected user", func(t *testing.T) {
		sut, tearDown := c.NewSecretRepository()
		t.Cleanup(tearDown)
		userID := uuid.New()
		ctx := context.Background()
		secret := &model.Secret{}
		s1, err := sut.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		secret = &model.Secret{}
		s2, err := sut.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		user2ID := uuid.New()
		secret = &model.Secret{}
		_, err = sut.AddSecret(ctx, secret, user2ID)
		require.NoError(t, err)
		want := []*model.Secret{
			{ID: s1.ID},
			{ID: s2.ID},
		}

		got, err := sut.ListSecrets(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
