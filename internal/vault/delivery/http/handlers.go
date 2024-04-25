package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/vault"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

type ListDataResponse struct {
	List []Secret `json:"list,omitempty"`
}

type Secret struct {
}

type VaultHandlers struct {
	authConfig config.JWTAuthConfig
}

func NewVaultHandlers(authConfig config.JWTAuthConfig) vault.VaultHandlers {
	return &VaultHandlers{
		authConfig: authConfig,
	}
}

func (h *VaultHandlers) ListSecrets() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		_ = claims

		resp := ListDataResponse{}
		content, _ := json.Marshal(resp)

		w.Header().Set(contentTypeHeader, applicationJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(content)
	})
}
