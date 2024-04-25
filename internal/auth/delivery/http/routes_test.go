package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/auth/service"
	"github.com/nestjam/goph-keeper/internal/config"
)

func TestMapAuthRoutes(t *testing.T) {
	const (
		email        = "user123@email.com"
		password     = "password"
		registerPath = "/register"
	)

	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

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
	t.Run("regiser user with plain text content type", func(t *testing.T) {
		service := &authServiceMock{}
		handlers := NewAuthHandlers(service, config)
		sut := chi.NewRouter()

		MapAuthRoutes(sut, handlers)
		r := newPlainTextRequest(t, registerPath)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	})
}

func newPlainTextRequest(t *testing.T, target string) *http.Request {
	t.Helper()

	body := "plaint text"
	r := httptest.NewRequest(http.MethodPost, target, strings.NewReader(body))
	r.Header.Set(contentTypeHeader, "text/plain")
	return r
}
