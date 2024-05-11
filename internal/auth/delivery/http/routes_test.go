package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/auth/service"
	"github.com/nestjam/goph-keeper/internal/config"
)

func TestMapAuthRoutes(t *testing.T) {
	const (
		email        = "user123@email.com"
		password     = "password"
		registerPath = "/register"
		loginPath    = "/login"
	)

	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

	t.Run("register", func(t *testing.T) {
		t.Run("regiser user", func(t *testing.T) {
			repo := inmemory.NewUserRepository()
			service := service.NewAuthService(repo)
			handlers := NewAuthHandlers(service, config)
			sut := chi.NewRouter()

			MapAuthRoutes(sut, handlers)
			r := newRegisterUserRequest(t, registerPath, email, password)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, http.StatusCreated, w.Code)
		})
		t.Run("request to regiser user with plain text content type", func(t *testing.T) {
			service := &authServiceMock{}
			handlers := NewAuthHandlers(service, config)
			sut := chi.NewRouter()

			MapAuthRoutes(sut, handlers)
			r := newPlainTextRequest(t, registerPath)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
		})
	})
	t.Run("login", func(t *testing.T) {
		t.Run("login user", func(t *testing.T) {
			ctx := context.Background()
			repo := inmemory.NewUserRepository()
			registerUser(t, ctx, email, password, repo)
			service := service.NewAuthService(repo)
			handlers := NewAuthHandlers(service, config)
			sut := chi.NewRouter()

			MapAuthRoutes(sut, handlers)
			r := newLoginUserRequest(t, loginPath, email, password)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}

func registerUser(t *testing.T, ctx context.Context, email, password string, repo auth.UserRepository) {
	t.Helper()

	user := &model.User{Email: email, Password: password}
	err := user.HashPassword()
	require.NoError(t, err)
	_, err = repo.Register(ctx, user)
	require.NoError(t, err)
}

func newPlainTextRequest(t *testing.T, target string) *http.Request {
	t.Helper()

	body := "plaint text"
	r := httptest.NewRequest(http.MethodPost, target, strings.NewReader(body))
	r.Header.Set(contentTypeHeader, "text/plain")
	return r
}
