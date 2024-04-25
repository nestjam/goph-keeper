package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
)

func TestMapVaultRoutes(t *testing.T) {
	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

	t.Run("list secrets", func(t *testing.T) {
		handlers := &vaultHandlersMock{}
		sut := chi.NewRouter()

		MapVaultRoutes(sut, handlers, config)
		r := newListSecretsRequest(t, "/list")
		setAuthCookie(t, r, config)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, 1, handlers.listSecretsCalls)
	})
	t.Run("user is not authenticated to list secrets", func(t *testing.T) {
		handlers := &vaultHandlersMock{}
		sut := chi.NewRouter()

		MapVaultRoutes(sut, handlers, config)
		r := newListSecretsRequest(t, "/list")
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, 0, handlers.listSecretsCalls)
	})
}

func setAuthCookie(t *testing.T, r *http.Request, config config.JWTAuthConfig) {
	t.Helper()

	baker := utils.NewAuthCookieBaker(config)
	cookie, err := baker.BakeCookie(uuid.New())
	require.NoError(t, err)
	r.AddCookie(cookie)
}
