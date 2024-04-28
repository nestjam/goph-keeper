package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
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

		resp := createListSecretsResponse(secrets)
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

		secret, err := getSecret(r)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		addedSecret, err := h.service.AddSecret(ctx, secret, userID)
		if err != nil {
			writeInternalServerError(w)
			return
		}

		resp := createAddSecretResponse(addedSecret)
		err = writeJSON(w, http.StatusCreated, resp)
		if err != nil {
			writeInternalServerError(w)
			return
		}
	})
}

func getSecret(r *http.Request) (*model.Secret, error) {
	const op = "get secret"

	var req AddSecretRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	secret := &model.Secret{
		Data: req.Secret.Data,
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

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func createListSecretsResponse(secrets []*model.Secret) *ListSecretsResponse {
	resp := &ListSecretsResponse{
		List: make([]Secret, len(secrets)),
	}

	for i := 0; i < len(secrets); i++ {
		s := secrets[i]
		resp.List[i] = Secret{ID: s.ID.String()}
	}

	return resp
}

func createAddSecretResponse(secret *model.Secret) *AddSecretResponse {
	return &AddSecretResponse{
		Secret: Secret{
			ID: secret.ID.String(),
		},
	}
}
