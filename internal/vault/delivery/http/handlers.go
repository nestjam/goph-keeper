package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
	secretParam       = "secret"
)

type VaultHandlers struct {
	service    vault.VaultService
	authConfig config.JWTAuthConfig
}

func NewVaultHandlers(service vault.VaultService, authConfig config.JWTAuthConfig) vault.VaultHandlers {
	return &VaultHandlers{
		service:    service,
		authConfig: authConfig,
	}
}

func (h *VaultHandlers) ListSecrets() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := utils.UserFromContext(ctx)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		secrets, err := h.service.ListSecrets(ctx, userID)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		resp := newListSecretsResponse(secrets)
		err = writeJSON(w, http.StatusOK, resp)
		if err != nil {
			writeInternalServerError(w)
			return
		}
	})
}

func (h *VaultHandlers) AddSecret() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := utils.UserFromContext(ctx)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		secret, err := secretFromAddRequest(r)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		secretID, err := h.service.AddSecret(ctx, secret, userID)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		resp := newAddSecretResponse(secretID)
		err = writeJSON(w, http.StatusCreated, resp)
		if err != nil {
			writeInternalServerError(w)
			return
		}
	})
}

func (h *VaultHandlers) UpdateSecret() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, secretParam)
		secretID, err := uuid.Parse(key)
		if err != nil {
			writeBadRequest(w)
			return
		}

		ctx := r.Context()
		userID, err := utils.UserFromContext(ctx)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		secret, err := secretFromUpdateRequest(r)
		if err != nil {
			writeInternalServerError(w)
			return
		}
		secret.ID = secretID

		err = h.service.UpdateSecret(ctx, secret, userID)
		if errors.Is(err, vault.ErrSecretNotFound) {
			writeNotFound(w)
			return
		}
		if err != nil {
			writeInternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (h *VaultHandlers) GetSecret() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, secretParam)
		secretID, err := uuid.Parse(key)
		if err != nil {
			writeBadRequest(w)
			return
		}

		ctx := r.Context()
		userID, err := utils.UserFromContext(ctx)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		secret, err := h.service.GetSecret(ctx, secretID, userID)
		if errors.Is(err, vault.ErrSecretNotFound) {
			writeNotFound(w)
			return
		}
		if err != nil {
			writeInternalServerError(w)
			return
		}

		resp := newGetSecretResponse(secret)
		err = writeJSON(w, http.StatusOK, resp)
		if err != nil {
			writeInternalServerError(w)
			return
		}
	})
}

func (h *VaultHandlers) DeleteSecret() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, secretParam)
		secretID, err := uuid.Parse(key)
		if err != nil {
			writeBadRequest(w)
			return
		}

		ctx := r.Context()
		userID, err := utils.UserFromContext(ctx)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		err = h.service.DeleteSecret(ctx, secretID, userID)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func newGetSecretResponse(secret *model.Secret) GetSecretResponse {
	return GetSecretResponse{
		Secret: Secret{
			ID:   secret.ID.String(),
			Name: secret.Name,
			Data: string(secret.Data),
		},
	}
}

func secretFromAddRequest(r *http.Request) (*model.Secret, error) {
	const op = "get secret"

	var req AddSecretRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	secret := &model.Secret{
		Name: req.Secret.Name,
		Data: []byte(req.Secret.Data),
	}
	return secret, nil
}

func secretFromUpdateRequest(r *http.Request) (*model.Secret, error) {
	const op = "update secret"

	var req UpdateSecretRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	secret := &model.Secret{
		Name: req.Secret.Name,
		Data: []byte(req.Secret.Data),
	}
	return secret, nil
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) error {
	const op = "write json"

	content, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, op)
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(statusCode)
	_, _ = w.Write(content)
	return nil
}

func writeBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func writeNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func newListSecretsResponse(secrets []*model.Secret) *ListSecretsResponse {
	resp := &ListSecretsResponse{
		List: make([]Secret, len(secrets)),
	}

	for i := 0; i < len(secrets); i++ {
		s := secrets[i]
		resp.List[i] = Secret{
			ID:   s.ID.String(),
			Name: s.Name,
		}
	}

	return resp
}

func newAddSecretResponse(secretID uuid.UUID) AddSecretResponse {
	return AddSecretResponse{
		Secret: Secret{
			ID: secretID.String(),
		},
	}
}
