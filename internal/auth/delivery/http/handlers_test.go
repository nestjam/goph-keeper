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
	"github.com/google/uuid"
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

		assertAuthToken(t, w, config, user.ID)
	})
	t.Run("register request contains invalid json", func(t *testing.T) {
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
		service.RegisterFunc = func(ctx context.Context, user *model.User) (uuid.UUID, error) {
			return uuid.Nil, errors.New("failed to register")
		}
		sut := NewAuthHandlers(service, config)
		r := newRegisterUserRequest(t, "/", email, password)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLogin(t *testing.T) {
	const (
		email    = "user@email.com"
		password = "1234"
	)

	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

	t.Run("login registered user", func(t *testing.T) {
		ctx := context.Background()
		user := &model.User{Email: email, Password: password}
		err := user.HashPassword()
		require.NoError(t, err)
		repo := inmemory.NewUserRepository()
		user.ID, err = repo.Register(ctx, user)
		require.NoError(t, err)
		service := service.NewAuthService(repo)
		sut := NewAuthHandlers(service, config)
		r := newLoginUserRequest(t, "/", email, password)
		w := httptest.NewRecorder()

		sut.Login().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assertAuthToken(t, w, config, user.ID)
	})
	t.Run("login request contains invalid json", func(t *testing.T) {
		repo := inmemory.NewUserRepository()
		service := service.NewAuthService(repo)
		sut := NewAuthHandlers(service, config)
		r := newLoginUserInvalidRequest(t)
		w := httptest.NewRecorder()

		sut.Login().ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("login failed", func(t *testing.T) {
		service := &authServiceMock{}
		service.LoginFunc = func(ctx context.Context, user *model.User) (uuid.UUID, error) {
			return uuid.Nil, errors.New("failed to login")
		}
		sut := NewAuthHandlers(service, config)
		r := newLoginUserRequest(t, "/", email, password)
		w := httptest.NewRecorder()

		sut.Login().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func newRegisterUserInvalidRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{{register}"))
}

func newLoginUserInvalidRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{{login}"))
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

func newLoginUserRequest(t *testing.T, path, email, password string) *http.Request {
	t.Helper()

	data := LoginUserRequest{
		Email:    email,
		Password: password,
	}
	body, err := json.Marshal(data)
	require.NoError(t, err)
	r := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	r.Header.Set(contentTypeHeader, applicationJSON)
	return r
}

func assertAuthToken(t *testing.T, w *httptest.ResponseRecorder, cfg config.JWTAuthConfig, user uuid.UUID) {
	t.Helper()

	r := w.Result()
	defer func() { _ = r.Body.Close() }()
	cookies := r.Cookies()
	assert.NotEmpty(t, cookies)
	jwtCookie := cookies[0]
	assert.Equal(t, utils.JWTCookieName, jwtCookie.Name)
	wantMaxAge := int(cfg.TokenExpiryIn / time.Second)
	assert.Equal(t, wantMaxAge, jwtCookie.MaxAge)
	assert.Equal(t, true, jwtCookie.HttpOnly)

	jwtAuth := jwtauth.New(utils.JWTAlg, []byte(cfg.SignKey), nil)
	token, err := jwtAuth.Decode(jwtCookie.Value)
	require.NoError(t, err)
	id, ok := token.Get(utils.UserIDClaim)
	require.True(t, ok)
	assert.Equal(t, user.String(), id.(string))
}
