package model

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("hash password", func(t *testing.T) {
		const password = "1234"
		sut := &User{Password: password}

		err := sut.HashPassword()

		require.NoError(t, err)
		assertEqualPasswords(t, password, sut.Password)
	})
	t.Run("empty password", func(t *testing.T) {
		const password = ""
		sut := &User{Password: password}

		err := sut.HashPassword()

		require.NoError(t, err)
		assertEqualPasswords(t, password, sut.Password)
	})
	t.Run("password is longer than 72 symbols", func(t *testing.T) {
		sut := &User{Password: strings.Repeat("0", PasswordMaxLengthInBytes+1)}

		err := sut.HashPassword()

		require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	})
}

func TestComparePassword(t *testing.T) {
	t.Run("passwords are equal", func(t *testing.T) {
		const password = "1234"
		sut := &User{Password: password}
		err := sut.HashPassword()
		require.NoError(t, err)

		got := sut.ComparePassword(password)

		assert.True(t, got)
	})
	t.Run("passwords are not equal", func(t *testing.T) {
		const password = "1234"
		sut := &User{Password: password}
		err := sut.HashPassword()
		require.NoError(t, err)
		const p = "4321"

		got := sut.ComparePassword(p)

		assert.False(t, got)
	})
}

func TestCopy(t *testing.T) {
	sut := &User{
		ID:       uuid.New(),
		Email:    "user@email.com",
		Password: "1234",
	}

	got := sut.Copy()
	assert.True(t, sut != got)
	assert.Equal(t, sut, got)
}

func assertEqualPasswords(t *testing.T, password, hashedPassword string) {
	t.Helper()

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}
