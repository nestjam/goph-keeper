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
	const (
		listPath = "/list"
	)

	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

	t.Run("list secrets", func(t *testing.T) {
		handlers := &vaultHandlersMock{}
		sut := chi.NewRouter()

		MapVaultRoutes(sut, handlers, config)
		r := newListSecretsRequest(t, listPath)
		want := uuid.New()
		setAuthCookie(t, r, config, want)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, 1, handlers.calls)
		got, ok := handlers.claims[utils.UserIDClaim]
		require.True(t, ok)
		assert.Equal(t, want.String(), got)
	})
	t.Run("user is not authenticated to list secrets", func(t *testing.T) {
		handlers := &vaultHandlersMock{}
		sut := chi.NewRouter()

		MapVaultRoutes(sut, handlers, config)
		r := newListSecretsRequest(t, listPath)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, 0, handlers.calls)
	})
}

func setAuthCookie(t *testing.T, r *http.Request, cfg config.JWTAuthConfig, id uuid.UUID) {
	t.Helper()

	baker := utils.NewAuthCookieBaker(cfg)
	cookie, err := baker.BakeCookie(id)
	require.NoError(t, err)
	r.AddCookie(cookie)
}
