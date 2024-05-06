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
		secret := &model.Secret{
			Data: []byte("data"),
		}
		userID := uuid.New()
		ctx := context.Background()

		got, err := sut.AddSecret(ctx, secret, userID)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, got.ID)
		assert.Equal(t, secret.Data, got.Data)
	})

	t.Run("update secret", func(t *testing.T) {
		sut, tearDown := c.NewSecretRepository()
		t.Cleanup(tearDown)
		secret := &model.Secret{}
		userID := uuid.New()
		ctx := context.Background()
		secret, err := sut.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		secret.Data = []byte("edited text")

		err = sut.UpdateSecret(ctx, secret, userID)

		require.NoError(t, err)
		got, _ := sut.GetSecret(ctx, secret.ID, userID)
		assert.Equal(t, secret.ID, got.ID)
		assert.Equal(t, secret.Data, got.Data)
	})
	t.Run("update secret that does not exist", func(t *testing.T) {
		sut, tearDown := c.NewSecretRepository()
		t.Cleanup(tearDown)
		secret := &model.Secret{ID: uuid.New()}
		userID := uuid.New()
		ctx := context.Background()

		err := sut.UpdateSecret(ctx, secret, userID)

		require.Error(t, err)
	})

	t.Run("list secrets", func(t *testing.T) {
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
			secret := &model.Secret{
				Data: []byte("data_"),
			}
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

		t.Run("get secret", func(t *testing.T) {
			t.Run("get user secret", func(t *testing.T) {
				sut, tearDown := c.NewSecretRepository()
				t.Cleanup(tearDown)
				secret := &model.Secret{}
				userID := uuid.New()
				ctx := context.Background()
				want, err := sut.AddSecret(ctx, secret, userID)
				require.NoError(t, err)

				got, err := sut.GetSecret(ctx, want.ID, userID)

				require.NoError(t, err)
				assert.Equal(t, want, got)
			})
			t.Run("user does not have the secret", func(t *testing.T) {
				sut, tearDown := c.NewSecretRepository()
				t.Cleanup(tearDown)
				secret := &model.Secret{}
				userID := uuid.New()
				ctx := context.Background()
				_, err := sut.AddSecret(ctx, secret, userID)
				require.NoError(t, err)
				anotherSecretID := uuid.New()

				_, err = sut.GetSecret(ctx, anotherSecretID, userID)

				require.ErrorIs(t, err, ErrSecretNotFound)
			})
			t.Run("user with id does not exist", func(t *testing.T) {
				sut, tearDown := c.NewSecretRepository()
				t.Cleanup(tearDown)
				userID := uuid.New()
				ctx := context.Background()
				secretID := uuid.New()

				_, err := sut.GetSecret(ctx, secretID, userID)

				require.ErrorIs(t, err, ErrSecretNotFound)
			})
		})

		t.Run("delete secret", func(t *testing.T) {
			t.Run("delete user secret", func(t *testing.T) {
				sut, tearDown := c.NewSecretRepository()
				t.Cleanup(tearDown)
				secret := &model.Secret{}
				userID := uuid.New()
				ctx := context.Background()
				want, err := sut.AddSecret(ctx, secret, userID)
				require.NoError(t, err)

				err = sut.DeleteSecret(ctx, want.ID, userID)

				require.NoError(t, err)
			})
			t.Run("user does not have the secret (on delete secret)", func(t *testing.T) {
				sut, tearDown := c.NewSecretRepository()
				t.Cleanup(tearDown)
				secret := &model.Secret{}
				userID := uuid.New()
				ctx := context.Background()
				_, err := sut.AddSecret(ctx, secret, userID)
				require.NoError(t, err)
				anotherSecretID := uuid.New()

				err = sut.DeleteSecret(ctx, anotherSecretID, userID)

				require.NoError(t, err)
			})
			t.Run("user with id does not exist (on delete secret)", func(t *testing.T) {
				sut, tearDown := c.NewSecretRepository()
				t.Cleanup(tearDown)
				userID := uuid.New()
				ctx := context.Background()
				secretID := uuid.New()

				err := sut.DeleteSecret(ctx, secretID, userID)

				require.NoError(t, err)
			})
		})
	})
}
