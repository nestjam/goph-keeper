package http

import (
	"encoding/json"
	"net/http"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/pkg/errors"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

type ListSecretsResponse struct {
	List []SecretInfo `json:"list,omitempty"`
}

type SecretInfo struct {
	ID string `json:"id"`
}

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
			setInternalServerError(w)
			return
		}

		secrets, err := h.service.ListSecrets(ctx, userID)
		if err != nil {
			setInternalServerError(w)
			return
		}

		resp := createListSecretsResponse(secrets)
		err = setOK(w, resp)
		if err != nil {
			setInternalServerError(w)
			return
		}
	})
}

func setOK(w http.ResponseWriter, v any) error {
	const op = "set OK"

	content, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, op)
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
	return nil
}

func setInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func createListSecretsResponse(secrets []*model.Secret) *ListSecretsResponse {
	resp := &ListSecretsResponse{
		List: make([]SecretInfo, len(secrets)),
	}

	for i := 0; i < len(secrets); i++ {
		s := secrets[i]
		resp.List[i] = SecretInfo{ID: s.ID.String()}
	}

	return resp
}
