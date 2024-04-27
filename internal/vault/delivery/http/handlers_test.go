package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/nestjam/goph-keeper/internal/vault/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/vault/service"
)

func TestListSecrets(t *testing.T) {
	config := getConfig()

	t.Run("empty list", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		service := service.NewVaultService(repo)
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		r := newListSecretsRequestWithUser(t, userID)
		w := httptest.NewRecorder()

		sut.ListSecrets().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assertContentType(t, applicationJSON, w)
		resp := getListSecretsResponse(t, w.Body)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.List)
	})
	t.Run("secrets", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		service := service.NewVaultService(repo)
		sut := NewVaultHandlers(service, config)
		ctx := context.Background()
		userID := uuid.New()
		secret := &model.Secret{}
		s, err := repo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		r := newListSecretsRequestWithUser(t, userID)
		w := httptest.NewRecorder()
		want := ListSecretsResponse{
			List: []Secret{
				{ID: s.ID.String()},
			},
		}

		sut.ListSecrets().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		got := getListSecretsResponse(t, w.Body)
		assert.Equal(t, want, got)
	})
	t.Run("user not found in context", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		service := service.NewVaultService(repo)
		sut := NewVaultHandlers(service, config)
		r := newListSecretsRequest(t, "/")
		r = addAuthError(t, r, errors.New("failed"))
		w := httptest.NewRecorder()

		sut.ListSecrets().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("failed to list secrets", func(t *testing.T) {
		service := &vaultServiceMock{
			ListSecretsFunc: func(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
				return nil, errors.New("failed")
			},
		}
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		r := newListSecretsRequestWithUser(t, userID)
		w := httptest.NewRecorder()

		sut.ListSecrets().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAddSecret(t *testing.T) {
	config := getConfig()

	t.Run("add secret", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		service := service.NewVaultService(repo)
		sut := NewVaultHandlers(service, config)
		payload := "sensitive data"
		secret := Secret{Payload: payload}
		userID := uuid.New()
		r := newAddSecretRequestWithUser(t, secret, userID)
		w := httptest.NewRecorder()

		sut.AddSecret().ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		ctx := context.Background()
		secrets, err := repo.ListSecrets(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(secrets))
		want := secrets[0]
		assert.Equal(t, payload, want.Payload)

		resp := getAddSecretResponse(t, w.Body)
		assert.Equal(t, want.ID.String(), resp.Secret.ID)
	})
	t.Run("user not found in context", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		service := service.NewVaultService(repo)
		sut := NewVaultHandlers(service, config)
		secret := Secret{}
		r := newAddSecretRequest(t, "/", secret)
		r = addAuthError(t, r, errors.New("failed"))
		w := httptest.NewRecorder()

		sut.AddSecret().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("failed to add secret", func(t *testing.T) {
		service := &vaultServiceMock{
			AddSecretFunc: func(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error) {
				return nil, errors.New("failed")
			},
		}
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		secret := Secret{}
		r := newAddSecretRequestWithUser(t, secret, userID)
		w := httptest.NewRecorder()

		sut.AddSecret().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("invalid json", func(t *testing.T) {
		repo := inmemory.NewSecretRepository()
		service := service.NewVaultService(repo)
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		r := newInvalidAddSecretRequestWithUser(t, userID)
		w := httptest.NewRecorder()

		sut.AddSecret().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func getConfig() config.JWTAuthConfig {
	return config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}
}

func newAddSecretRequestWithUser(t *testing.T, secret Secret, userID uuid.UUID) *http.Request {
	t.Helper()

	r := newAddSecretRequest(t, "/", secret)
	r = addAuthToken(t, r, userID)
	return r
}

func newInvalidAddSecretRequestWithUser(t *testing.T, userID uuid.UUID) *http.Request {
	t.Helper()

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{{invalid json]}"))
	r = addAuthToken(t, r, userID)
	return r
}

func newListSecretsRequest(t *testing.T, path string) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodGet, path, nil)
}

func newListSecretsRequestWithUser(t *testing.T, userID uuid.UUID) *http.Request {
	t.Helper()

	r := newListSecretsRequest(t, "/")
	r = addAuthToken(t, r, userID)
	return r
}

func getListSecretsResponse(t *testing.T, r io.Reader) ListSecretsResponse {
	t.Helper()

	var resp ListSecretsResponse
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&resp)
	require.NoError(t, err)
	return resp
}

func getAddSecretResponse(t *testing.T, r io.Reader) AddSecretResponse {
	t.Helper()

	var resp AddSecretResponse
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&resp)
	require.NoError(t, err)
	return resp
}

func addAuthToken(t *testing.T, r *http.Request, userID uuid.UUID) *http.Request {
	t.Helper()

	token := jwt.New()
	err := token.Set(utils.UserIDClaim, userID.String())
	require.NoError(t, err)

	ctx := r.Context()
	ctx = context.WithValue(ctx, jwtauth.TokenCtxKey, token)
	return r.WithContext(ctx)
}

func addAuthError(t *testing.T, r *http.Request, err error) *http.Request {
	t.Helper()

	ctx := r.Context()
	ctx = context.WithValue(ctx, jwtauth.ErrorCtxKey, err)
	return r.WithContext(ctx)
}

func assertContentType(t *testing.T, want string, r *httptest.ResponseRecorder) {
	t.Helper()

	assert.Equal(t, want, r.Header().Get(contentTypeHeader))
}
