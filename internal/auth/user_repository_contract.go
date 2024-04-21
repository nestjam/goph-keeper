package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserRepositoryContract struct {
	NewUserRepository func() (UserRepository, func())
}

func (c UserRepositoryContract) Test(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		t.Run("register new user", func(t *testing.T) {
			sut, tearDown := c.NewUserRepository()
			t.Cleanup(tearDown)

			user := &model.User{
				Email:    "user@email.com",
				Password: "123",
			}
			got, err := sut.Register(user)

			require.NoError(t, err)
			assert.Equal(t, user.Email, got.Email)
			assert.Equal(t, user.Password, got.Password)
			assert.NotEqual(t, uuid.UUID{}, got.ID)
		})
		t.Run("register user with email that has already been registered", func(t *testing.T) {
			sut, tearDown := c.NewUserRepository()
			t.Cleanup(tearDown)
			user := &model.User{
				Email:    "user@email.com",
				Password: "123",
			}
			_, err := sut.Register(user)
			require.NoError(t, err)

			_, err = sut.Register(user)
			require.ErrorIs(t, err, ErrUserWithEmailIsRegistered)
		})
	})

	t.Run("find user by email", func(t *testing.T) {
		t.Run("find existing user", func(t *testing.T) {
			sut, tearDown := c.NewUserRepository()
			t.Cleanup(tearDown)
			user := &model.User{
				Email: "user@email.com",
			}
			want, err := sut.Register(user)
			require.NoError(t, err)

			got, err := sut.FindByEmail(user.Email)
			require.NoError(t, err)
			assert.Equal(t, want.ID, got.ID)
			assert.Equal(t, user.Email, got.Email)
		})
		t.Run("find user that does not exist", func(t *testing.T) {
			sut, tearDown := c.NewUserRepository()
			t.Cleanup(tearDown)
			user := &model.User{
				Email: "user@email.com",
			}

			_, err := sut.FindByEmail(user.Email)
			require.ErrorIs(t, err, ErrUserIsNotRegisteredAtEmail)
		})
	})
}
