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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/nestjam/goph-keeper/internal/vault/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/vault/service"
)

func TestListSecrets(t *testing.T) {
	config := newConfig()
	rootKey := randomMasterKey(t)

	t.Run("empty list", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		r := newListSecretsRequestWithUser(t, userID)
		w := httptest.NewRecorder()

		sut.ListSecrets().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assertContentType(t, applicationJSON, w)
		got := listSecretsFromResponse(t, w.Body)
		assert.Empty(t, got)
	})
	t.Run("secrets", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		ctx := context.Background()
		userID := uuid.New()
		secret := &model.Secret{}
		s, err := secretRepo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		r := newListSecretsRequestWithUser(t, userID)
		w := httptest.NewRecorder()
		want := []Secret{
			{ID: s.ID.String()},
		}

		sut.ListSecrets().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		got := listSecretsFromResponse(t, w.Body)
		assert.Equal(t, want, got)
	})
	t.Run("user not found in context", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
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
	config := newConfig()
	rootKey := randomMasterKey(t)

	t.Run("add secret", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		const data = "sensitive data"
		secret := Secret{Data: data}
		userID := uuid.New()
		r := newAddSecretRequestWithUser(t, secret, userID)
		w := httptest.NewRecorder()

		sut.AddSecret().ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		ctx := context.Background()
		secrets, err := secretRepo.ListSecrets(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(secrets))
		want := secrets[0]

		resp := getAddSecretResponse(t, w.Body)
		assert.Equal(t, want.ID.String(), resp.Secret.ID)
	})
	t.Run("user not found in context", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
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
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		r := newInvalidAddSecretRequestWithUser(t, userID)
		w := httptest.NewRecorder()

		sut.AddSecret().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateSecret(t *testing.T) {
	config := newConfig()
	rootKey := randomMasterKey(t)

	t.Run("update secret", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		ctx := context.Background()
		s := &model.Secret{}
		userID := uuid.New()
		tt, err := service.AddSecret(ctx, s, userID)
		s.ID = tt.ID
		require.NoError(t, err)
		const wantData = "edited text"
		secret := Secret{ID: s.ID.String(), Data: wantData}
		r := newUpdateSecretRequestWithUser(t, secret, userID)
		w := httptest.NewRecorder()
		sut := NewVaultHandlers(service, config)

		updateSecret(sut, w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		got, err := service.GetSecret(ctx, s.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, []byte(wantData), got.Data)
	})
	t.Run("updating secret not found", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		userID := uuid.New()
		const wantData = "edited text"
		secret := Secret{ID: uuid.NewString(), Data: wantData}
		r := newUpdateSecretRequestWithUser(t, secret, userID)
		w := httptest.NewRecorder()
		sut := NewVaultHandlers(service, config)

		updateSecret(sut, w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("user not found in context", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		secret := Secret{ID: uuid.NewString()}
		r := newUpdateSecretRequest(t, "", secret)
		r = addAuthError(t, r, errors.New("failed"))
		w := httptest.NewRecorder()

		updateSecret(sut, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("invalid secret id", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		r := newInvalidIDUpdateSecretRequest(t)
		w := httptest.NewRecorder()

		updateSecret(sut, w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid json", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		secretID := uuid.NewString()
		r := newInvalidUpdateSecretRequestWithUser(t, secretID, userID)
		w := httptest.NewRecorder()

		updateSecret(sut, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("failed to update secret", func(t *testing.T) {
		service := &vaultServiceMock{
			UpdateSecretFunc: func(ctx context.Context, secret *model.Secret, userID uuid.UUID) error {
				return errors.New("failed")
			},
		}
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		secret := Secret{ID: uuid.NewString()}
		r := newUpdateSecretRequestWithUser(t, secret, userID)
		w := httptest.NewRecorder()

		updateSecret(sut, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetSecret(t *testing.T) {
	config := newConfig()
	rootKey := randomMasterKey(t)

	t.Run("get secret", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		ctx := context.Background()
		userID := uuid.New()
		const data = "data"
		secret := &model.Secret{Data: []byte(data)}
		added, err := service.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		r := newGetSecretRequestWithUser(t, added.ID, userID)
		w := httptest.NewRecorder()

		getSecret(sut, w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp := secretFromResponse(t, w.Body)
		assert.Equal(t, added.ID.String(), resp.ID)
		assert.Equal(t, data, resp.Data)
	})
	t.Run("secret not found", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		secretID := uuid.New()
		r := newGetSecretRequestWithUser(t, secretID, userID)
		w := httptest.NewRecorder()

		getSecret(sut, w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("invalid secret id", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		r := newInvalidIDGetSecretRequest(t)
		w := httptest.NewRecorder()

		getSecret(sut, w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("user not found in context", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		ctx := context.Background()
		userID := uuid.New()
		secret, err := secretRepo.AddSecret(ctx, &model.Secret{}, userID)
		require.NoError(t, err)
		r := newGetSecretRequest(t, "", secret.ID)
		r = addAuthError(t, r, errors.New("failed"))
		w := httptest.NewRecorder()

		getSecret(sut, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("failed to get secret", func(t *testing.T) {
		service := &vaultServiceMock{
			GetSecretFunc: func(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
				return nil, errors.New("failed")
			},
		}
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		secretID := uuid.New()
		r := newGetSecretRequestWithUser(t, secretID, userID)
		w := httptest.NewRecorder()

		getSecret(sut, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeleteSecret(t *testing.T) {
	config := newConfig()
	rootKey := randomMasterKey(t)

	t.Run("delete secret", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		ctx := context.Background()
		userID := uuid.New()
		secret := &model.Secret{}
		want, err := secretRepo.AddSecret(ctx, secret, userID)
		require.NoError(t, err)
		r := newDeleteSecretRequestWithUser(t, want.ID, userID)
		w := httptest.NewRecorder()

		deleteSecret(sut, w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("invalid secret id", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		r := newInvalidIDDeleteSecretRequest(t)
		w := httptest.NewRecorder()

		deleteSecret(sut, w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("user not found in context", func(t *testing.T) {
		keyRepo := inmemory.NewDataKeyRepository()
		secretRepo := inmemory.NewSecretRepository()
		service := service.NewVaultService(secretRepo, keyRepo, rootKey)
		sut := NewVaultHandlers(service, config)
		secretID := uuid.New()
		r := newDeleteSecretRequest(t, "", secretID)
		r = addAuthError(t, r, errors.New("failed"))
		w := httptest.NewRecorder()

		deleteSecret(sut, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("failed to delete secret", func(t *testing.T) {
		service := &vaultServiceMock{
			DeleteSecretFunc: func(ctx context.Context, secretID, userID uuid.UUID) error {
				return errors.New("failed")
			},
		}
		sut := NewVaultHandlers(service, config)
		userID := uuid.New()
		secretID := uuid.New()
		r := newDeleteSecretRequestWithUser(t, secretID, userID)
		w := httptest.NewRecorder()

		deleteSecret(sut, w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func updateSecret(sut vault.VaultHandlers, w *httptest.ResponseRecorder, r *http.Request) {
	router := chi.NewRouter()
	router.Patch("/{secret}", sut.UpdateSecret())
	router.ServeHTTP(w, r)
}

func deleteSecret(sut vault.VaultHandlers, w *httptest.ResponseRecorder, r *http.Request) {
	router := chi.NewRouter()
	router.Delete("/{secret}", sut.DeleteSecret())
	router.ServeHTTP(w, r)
}

func getSecret(sut vault.VaultHandlers, w *httptest.ResponseRecorder, r *http.Request) {
	router := chi.NewRouter()
	router.Get("/{secret}", sut.GetSecret())
	router.ServeHTTP(w, r)
}

func secretFromResponse(t *testing.T, r io.Reader) Secret {
	t.Helper()

	var resp GetSecretResponse
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&resp)
	require.NoError(t, err)
	return resp.Secret
}

func newDeleteSecretRequestWithUser(t *testing.T, secretID, userID uuid.UUID) *http.Request {
	t.Helper()

	r := newDeleteSecretRequest(t, "", secretID)
	r = addAuthToken(t, r, userID)
	return r
}

func newGetSecretRequestWithUser(t *testing.T, secretID, userID uuid.UUID) *http.Request {
	t.Helper()

	r := newGetSecretRequest(t, "", secretID)
	r = addAuthToken(t, r, userID)
	return r
}

func newInvalidIDUpdateSecretRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodPatch, "/foo", nil)
}

func newInvalidIDDeleteSecretRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodDelete, "/xyz", nil)
}

func newInvalidIDGetSecretRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodGet, "/abc", nil)
}

func newConfig() config.JWTAuthConfig {
	return config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}
}

func newUpdateSecretRequestWithUser(t *testing.T, secret Secret, userID uuid.UUID) *http.Request {
	t.Helper()

	r := newUpdateSecretRequest(t, "", secret)
	r = addAuthToken(t, r, userID)
	return r
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

func newInvalidUpdateSecretRequestWithUser(t *testing.T, secretID string, userID uuid.UUID) *http.Request {
	t.Helper()

	r := httptest.NewRequest(http.MethodPatch, "/"+secretID, strings.NewReader("{{invalid json]}"))
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

func listSecretsFromResponse(t *testing.T, r io.Reader) []Secret {
	t.Helper()

	var resp ListSecretsResponse
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&resp)
	require.NoError(t, err)
	return resp.List
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

func randomMasterKey(t *testing.T) *model.MasterKey {
	t.Helper()

	key, err := utils.GenerateRandomAES256Key()
	require.NoError(t, err)
	return model.NewMasterKey(key)
}
