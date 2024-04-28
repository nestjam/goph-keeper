package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
		secretsPath = "/secrets"
	)

	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

	t.Run("get secrets", func(t *testing.T) {
		t.Run("list secrets", func(t *testing.T) {
			spy := &vaultHandlersSpy{}
			sut := chi.NewRouter()

			MapVaultRoutes(sut, spy, config)
			r := newListSecretsRequest(t, secretsPath)
			want := uuid.New()
			setAuthCookie(t, r, config, want)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, 1, spy.listSecretsCallsCount)
			assertUserIDFromToken(t, want, spy)
		})
		t.Run("user is not authenticated to list secrets", func(t *testing.T) {
			handlers := &vaultHandlersSpy{}
			sut := chi.NewRouter()

			MapVaultRoutes(sut, handlers, config)
			r := newListSecretsRequest(t, secretsPath)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Equal(t, 0, handlers.listSecretsCallsCount)
		})
	})

	t.Run("post secret", func(t *testing.T) {
		t.Run("add secret", func(t *testing.T) {
			spy := &vaultHandlersSpy{}
			sut := chi.NewRouter()

			MapVaultRoutes(sut, spy, config)
			secret := Secret{}
			r := newAddSecretRequest(t, secretsPath, secret)
			want := uuid.New()
			setAuthCookie(t, r, config, want)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, 1, spy.addSecretCallsCount)
			assertUserIDFromToken(t, want, spy)
		})
		t.Run("request to add secret with plain text content type", func(t *testing.T) {
			spy := &vaultHandlersSpy{}
			sut := chi.NewRouter()
			MapVaultRoutes(sut, spy, config)
			r := newPlainTextRequest(t, secretsPath)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
		})
	})

	t.Run("get secret", func(t *testing.T) {
		spy := &vaultHandlersSpy{}
		sut := chi.NewRouter()

		MapVaultRoutes(sut, spy, config)
		secretID := uuid.New()
		r := newGetSecretRequest(t, secretsPath, secretID)
		userID := uuid.New()
		setAuthCookie(t, r, config, userID)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, 1, spy.getSecretCallsCount)
	})
}

func newGetSecretRequest(t *testing.T, path string, secretID uuid.UUID) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodGet, path+"/"+secretID.String(), nil)
}

func newPlainTextRequest(t *testing.T, target string) *http.Request {
	t.Helper()

	body := "plaint text"
	r := httptest.NewRequest(http.MethodPost, target, strings.NewReader(body))
	r.Header.Set(contentTypeHeader, "text/plain")
	return r
}

func assertUserIDFromToken(t *testing.T, userID uuid.UUID, spy *vaultHandlersSpy) {
	t.Helper()

	got, ok := spy.claims[utils.UserIDClaim]
	require.True(t, ok)
	assert.Equal(t, userID.String(), got)
}

func setAuthCookie(t *testing.T, r *http.Request, cfg config.JWTAuthConfig, id uuid.UUID) {
	t.Helper()

	baker := utils.NewAuthCookieBaker(cfg)
	cookie, err := baker.BakeCookie(id)
	require.NoError(t, err)
	r.AddCookie(cookie)
}

func newAddSecretRequest(t *testing.T, path string, secret Secret) *http.Request {
	t.Helper()

	req := AddSecretRequest{Secret: secret}
	content, err := json.Marshal(req)
	require.NoError(t, err)
	r := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(content))
	r.Header.Set(contentTypeHeader, applicationJSON)
	return r
}
