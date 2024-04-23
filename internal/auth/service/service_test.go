package service

import (
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

		got, err := sut.Register(user)

		assert.Equal(t, user.Email, got.Email)
		require.NoError(t, err)
		assertEqualPasswords(t, user.Password, got.Password)

		foundUser, err := repo.FindByEmail(email)
		require.NoError(t, err)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, foundUser.Email, got.Email)
	})
	t.Run("password is too long", func(t *testing.T) {
		const email = "user@email.com"
		repo := inmemory.NewUserRepository()
		sut := NewAuthService(repo)
		user := &model.User{
			Email:    email,
			Password: strings.Repeat("0", model.PasswordMaxLengthInBytes+1),
		}

		_, err := sut.Register(user)

		require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	})
	t.Run("register new user with email that has already been registered", func(t *testing.T) {
		const (
			email    = "user@mail.com"
			password = "1234"
		)
		repo := inmemory.NewUserRepository()
		_, _ = repo.Register(&model.User{Email: email, Password: "psw"})
		sut := NewAuthService(repo)
		user := &model.User{Email: email, Password: password}

		_, err := sut.Register(user)

		require.ErrorIs(t, err, auth.ErrUserWithEmailIsRegistered)
	})
}

func assertEqualPasswords(t *testing.T, password, hashedPassword string) {
	t.Helper()

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}
