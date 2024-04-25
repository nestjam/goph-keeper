package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/auth/service"
	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
)

func TestRegister(t *testing.T) {
	const (
		email    = "user@email.com"
		password = "1234"
	)

	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

	t.Run("regiser new user", func(t *testing.T) {
		repo := inmemory.NewUserRepository()
		service := service.NewAuthService(repo)
		sut := NewAuthHandlers(service, config)
		r := newRegisterUserRequest(t, "/", email, password)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		ctx := context.Background()
		_, err := repo.FindByEmail(ctx, email)
		require.NoError(t, err)
	})
	t.Run("add jwt cookie on success registration", func(t *testing.T) {
		repo := inmemory.NewUserRepository()
		service := service.NewAuthService(repo)
		sut := NewAuthHandlers(service, config)
		r := newRegisterUserRequest(t, "/", email, password)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		ctx := context.Background()
		user, err := repo.FindByEmail(ctx, email)
		require.NoError(t, err)

		res := w.Result()
		defer func() { _ = res.Body.Close() }()
		cookies := res.Cookies()
		assert.NotEmpty(t, cookies)
		jwtCookie := cookies[0]
		assert.Equal(t, utils.JWTCookieName, jwtCookie.Name)
		wantMaxAge := int(config.TokenExpiryIn / time.Second)
		assert.Equal(t, wantMaxAge, jwtCookie.MaxAge)
		assert.Equal(t, true, jwtCookie.HttpOnly)

		jwtAuth := jwtauth.New(utils.JWTAlg, []byte(config.SignKey), nil)
		token, err := jwtAuth.Decode(jwtCookie.Value)
		require.NoError(t, err)
		id, ok := token.Get(utils.UserIDClaim)
		require.True(t, ok)
		assert.Equal(t, user.ID.String(), id.(string))
	})
	t.Run("json is invalid", func(t *testing.T) {
		repo := inmemory.NewUserRepository()
		service := service.NewAuthService(repo)
		sut := NewAuthHandlers(service, config)
		r := newRegisterUserInvalidRequest(t)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("register failed", func(t *testing.T) {
		service := &authServiceMock{}
		service.RegisterFunc = func(ctx context.Context, user *model.User) (*model.User, error) {
			return nil, errors.New("failed to register")
		}
		sut := NewAuthHandlers(service, config)
		r := newRegisterUserRequest(t, "/", email, password)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func newRegisterUserInvalidRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{{invalid}"))
}

func newRegisterUserRequest(t *testing.T, path, email, password string) *http.Request {
	t.Helper()

	data := RegisterUserRequest{
		Email:    email,
		Password: password,
	}
	body, err := json.Marshal(data)
	require.NoError(t, err)
	r := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	r.Header.Set(contentTypeHeader, applicationJSON)
	return r
}
