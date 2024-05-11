package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
)

func TestRegister(t *testing.T) {
	t.Run("register new user", func(t *testing.T) {
		const (
			email    = "user@mail.com"
			password = "1234"
		)
		repo := inmemory.NewUserRepository()
		sut := NewAuthService(repo)
		user := &model.User{Email: email, Password: password}
		ctx := context.Background()

		_, err := sut.Register(ctx, user)

		require.NoError(t, err)
		foundUser, err := repo.FindByEmail(ctx, email)
		require.NoError(t, err)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, foundUser.Email, user.Email)
	})
	t.Run("password is too long", func(t *testing.T) {
		const email = "user@email.com"
		repo := inmemory.NewUserRepository()
		sut := NewAuthService(repo)
		user := &model.User{
			Email:    email,
			Password: strings.Repeat("0", model.PasswordMaxLengthInBytes+1),
		}
		ctx := context.Background()

		_, err := sut.Register(ctx, user)

		require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	})
	t.Run("register new user with email that has already been registered", func(t *testing.T) {
		const (
			email    = "user@mail.com"
			password = "1234"
		)
		repo := inmemory.NewUserRepository()
		ctx := context.Background()
		_, _ = repo.Register(ctx, &model.User{Email: email, Password: "psw"})
		sut := NewAuthService(repo)
		user := &model.User{Email: email, Password: password}

		_, err := sut.Register(ctx, user)

		require.ErrorIs(t, err, auth.ErrUserWithEmailIsRegistered)
	})
}

func TestLogin(t *testing.T) {
	t.Run("login registered user", func(t *testing.T) {
		const (
			email    = "user@mail.com"
			password = "1234"
		)
		ctx := context.Background()
		repo := inmemory.NewUserRepository()
		want := &model.User{Email: email, Password: password}
		_ = want.HashPassword()
		var err error
		want.ID, err = repo.Register(ctx, want)
		require.NoError(t, err)
		user := &model.User{Email: email, Password: password}
		sut := NewAuthService(repo)

		got, err := sut.Login(ctx, user)

		require.NoError(t, err)
		assert.Equal(t, want.ID, got)
	})
	t.Run("login with wrong password", func(t *testing.T) {
		const (
			email    = "user@mail.com"
			password = "1234"
		)
		ctx := context.Background()
		repo := inmemory.NewUserRepository()
		want := &model.User{Email: email, Password: password}
		_ = want.HashPassword()
		_, err := repo.Register(ctx, want)
		require.NoError(t, err)
		const invalidPassword = "4321"
		user := &model.User{Email: email, Password: invalidPassword}
		sut := NewAuthService(repo)

		_, err = sut.Login(ctx, user)

		require.ErrorIs(t, err, auth.ErrInvalidPassword)
	})
	t.Run("user is not registered by email", func(t *testing.T) {
		const (
			email = "no@mail.com"
		)
		ctx := context.Background()
		repo := inmemory.NewUserRepository()
		sut := NewAuthService(repo)
		user := &model.User{Email: email}

		_, err := sut.Login(ctx, user)

		require.ErrorIs(t, err, auth.ErrUserIsNotRegistered)
	})
}
