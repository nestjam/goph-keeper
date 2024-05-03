package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/config"
)

func TestUserFromContext(t *testing.T) {
	t.Run("context contains jwt token with user id", func(t *testing.T) {
		ctx := context.Background()
		want := uuid.New()
		token := jwt.New()
		err := token.Set(UserIDClaim, want.String())
		require.NoError(t, err)
		ctx = context.WithValue(ctx, jwtauth.TokenCtxKey, token)

		got, err := UserFromContext(ctx)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("context contains jwt decoding error", func(t *testing.T) {
		ctx := context.Background()
		id := uuid.New()
		token := jwt.New()
		err := token.Set(UserIDClaim, id.String())
		require.NoError(t, err)
		ctx = context.WithValue(ctx, jwtauth.TokenCtxKey, token)
		want := errors.New("test")
		ctx = context.WithValue(ctx, jwtauth.ErrorCtxKey, want)

		_, got := UserFromContext(ctx)

		require.ErrorIs(t, got, want)
	})
	t.Run("jwt claims does not contain user id", func(t *testing.T) {
		ctx := context.Background()
		token := jwt.New()
		ctx = context.WithValue(ctx, jwtauth.TokenCtxKey, token)

		_, got := UserFromContext(ctx)

		require.ErrorIs(t, got, ErrUserIDNotFound)
	})
	t.Run("user id is not uuid", func(t *testing.T) {
		ctx := context.Background()
		token := jwt.New()
		err := token.Set(UserIDClaim, "user123")
		require.NoError(t, err)
		ctx = context.WithValue(ctx, jwtauth.TokenCtxKey, token)

		_, err = UserFromContext(ctx)

		require.Error(t, err)
	})
}

func TestBakeCookie(t *testing.T) {
	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}
	sut := NewAuthCookieBaker(config)
	userID := uuid.New()
	wantMaxAge := int(config.TokenExpiryIn / time.Second)

	cookie, err := sut.BakeCookie(userID)

	require.NoError(t, err)
	assert.Equal(t, true, cookie.HttpOnly)
	assert.Equal(t, JWTCookieName, cookie.Name)
	assert.Equal(t, wantMaxAge, cookie.MaxAge)
	assertAuthToken(t, userID, cookie.Value, config.SignKey)
}

func assertAuthToken(t *testing.T, want uuid.UUID, tkn, key string) {
	t.Helper()

	jwtAuth := jwtauth.New(JWTAlg, []byte(key), nil)
	token, err := jwtAuth.Decode(tkn)
	require.NoError(t, err)
	claims := token.PrivateClaims()
	value, _ := claims[UserIDClaim].(string)
	got, err := uuid.Parse(value)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
