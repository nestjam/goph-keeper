package vault

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type SecretTestData struct {
	Users uuid.UUIDs
	Keys  uuid.UUIDs
}

type SecretRepositoryContract struct {
	NewSecretRepository func() (SecretRepository, func(), SecretTestData)
}

func (c SecretRepositoryContract) Test(t *testing.T) {
	t.Run("add secret", func(t *testing.T) {
		sut, tearDown, td := c.NewSecretRepository()
		t.Cleanup(tearDown)
		secret := &model.Secret{
			Data:  []byte("123"),
			KeyID: td.Keys[0],
		}
		userID := td.Users[0]
		ctx := context.Background()

		got, err := sut.AddSecret(ctx, secret, userID)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, got)
	})

	t.Run("update secret", func(t *testing.T) {
		sut, tearDown, td := c.NewSecretRepository()
		t.Cleanup(tearDown)
		secret := &model.Secret{KeyID: td.Keys[0]}
		userID := td.Users[0]
		ctx := context.Background()
		var err error
		secret.ID, err = sut.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		secret.Data = []byte("edited text")
		secret.Name = "secret"

		err = sut.UpdateSecret(ctx, secret, userID)

		require.NoError(t, err)
		got, _ := sut.GetSecret(ctx, secret.ID, userID)
		assert.Equal(t, secret, got)
	})
	t.Run("update secret that does not exist", func(t *testing.T) {
		sut, tearDown, _ := c.NewSecretRepository()
		t.Cleanup(tearDown)
		secret := &model.Secret{ID: uuid.New()}
		userID := uuid.New()
		ctx := context.Background()

		err := sut.UpdateSecret(ctx, secret, userID)

		require.Error(t, err)
	})

	t.Run("list secrets", func(t *testing.T) {
		t.Run("user has no secrets", func(t *testing.T) {
			sut, tearDown, td := c.NewSecretRepository()
			t.Cleanup(tearDown)
			userID := td.Keys[0]
			ctx := context.Background()

			got, err := sut.ListSecrets(ctx, userID)

			require.NoError(t, err)
			assert.Empty(t, got)
		})
		t.Run("list user secrets without sensitive data", func(t *testing.T) {
			sut, tearDown, td := c.NewSecretRepository()
			t.Cleanup(tearDown)
			userID := td.Users[0]
			ctx := context.Background()
			secret := &model.Secret{
				Name:  "secret",
				Data:  []byte("data_"),
				KeyID: td.Keys[0],
			}
			var err error
			secret.ID, err = sut.AddSecret(ctx, secret, userID)
			require.NoError(t, err)
			want := []*model.Secret{
				{
					ID:   secret.ID,
					Name: secret.Name,
				},
			}

			got, err := sut.ListSecrets(ctx, userID)

			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
		t.Run("list secrets of selected user", func(t *testing.T) {
			sut, tearDown, td := c.NewSecretRepository()
			t.Cleanup(tearDown)
			userID := td.Users[0]
			ctx := context.Background()
			s1 := &model.Secret{KeyID: td.Keys[0]}
			var err error
			s1.ID, err = sut.AddSecret(ctx, s1, userID)
			require.NoError(t, err)
			s2 := &model.Secret{KeyID: td.Keys[0]}
			s2.ID, err = sut.AddSecret(ctx, s2, userID)
			require.NoError(t, err)
			user2ID := td.Users[1]
			s3 := &model.Secret{KeyID: td.Keys[0]}
			_, err = sut.AddSecret(ctx, s3, user2ID)
			require.NoError(t, err)
			want := []*model.Secret{
				{ID: s1.ID},
				{ID: s2.ID},
			}

			got, err := sut.ListSecrets(ctx, userID)

			require.NoError(t, err)
			assert.ElementsMatch(t, want, got)
		})

		t.Run("get secret", func(t *testing.T) {
			t.Run("get user secret", func(t *testing.T) {
				sut, tearDown, td := c.NewSecretRepository()
				t.Cleanup(tearDown)
				want := &model.Secret{
					KeyID: td.Keys[0],
					Data:  []byte("data"),
					Name:  "secret",
				}
				userID := td.Users[0]
				ctx := context.Background()
				var err error
				want.ID, err = sut.AddSecret(ctx, want, userID)
				require.NoError(t, err)

				got, err := sut.GetSecret(ctx, want.ID, userID)

				require.NoError(t, err)
				assert.Equal(t, want, got)
			})
			t.Run("user does not have the secret", func(t *testing.T) {
				sut, tearDown, td := c.NewSecretRepository()
				t.Cleanup(tearDown)
				secret := &model.Secret{KeyID: td.Keys[0]}
				userID := td.Users[0]
				ctx := context.Background()
				_, err := sut.AddSecret(ctx, secret, userID)
				require.NoError(t, err)
				anotherSecretID := uuid.New()

				_, err = sut.GetSecret(ctx, anotherSecretID, userID)

				require.ErrorIs(t, err, ErrSecretNotFound)
			})
			t.Run("user with id does not exist", func(t *testing.T) {
				sut, tearDown, _ := c.NewSecretRepository()
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
				sut, tearDown, td := c.NewSecretRepository()
				t.Cleanup(tearDown)
				userID := td.Users[0]
				ctx := context.Background()
				want := &model.Secret{KeyID: td.Keys[0]}
				var err error
				want.ID, err = sut.AddSecret(ctx, want, userID)
				require.NoError(t, err)

				err = sut.DeleteSecret(ctx, want.ID, userID)

				require.NoError(t, err)
			})
			t.Run("user does not have the secret (on delete secret)", func(t *testing.T) {
				sut, tearDown, td := c.NewSecretRepository()
				t.Cleanup(tearDown)
				secret := &model.Secret{KeyID: td.Keys[0]}
				userID := td.Users[0]
				ctx := context.Background()
				_, err := sut.AddSecret(ctx, secret, userID)
				require.NoError(t, err)
				anotherSecretID := uuid.New()

				err = sut.DeleteSecret(ctx, anotherSecretID, userID)

				require.NoError(t, err)
			})
			t.Run("user with id does not exist (on delete secret)", func(t *testing.T) {
				sut, tearDown, _ := c.NewSecretRepository()
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
