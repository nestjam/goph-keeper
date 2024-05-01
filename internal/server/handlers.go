package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	repoAuth "github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
	serviceAuth "github.com/nestjam/goph-keeper/internal/auth/service"
	"github.com/nestjam/goph-keeper/internal/config"
	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
	repoVault "github.com/nestjam/goph-keeper/internal/vault/repository/inmemory"
	serviceVault "github.com/nestjam/goph-keeper/internal/vault/service"
)

func (s *Server) mapHandlers() http.Handler {
	jwtAuthConfig := config.JWTAuthConfig{
		SignKey:       "supersecret",
		TokenExpiryIn: time.Hour,
	}
	authRepo := repoAuth.NewUserRepository()
	authService := serviceAuth.NewAuthService(authRepo)
	authHandlers := httpAuth.NewAuthHandlers(authService, jwtAuthConfig)

	secretRepo := repoVault.NewSecretRepository()
	keyRepo := repoVault.NewDataKeyRepository()
	vaultService := serviceVault.NewVaultService(secretRepo, keyRepo, s.rootKey)
	vaultHandlers := httpVault.NewVaultHandlers(vaultService, jwtAuthConfig)

	r := chi.NewRouter()
	httpAuth.MapAuthRoutes(r, authHandlers)
	httpVault.MapVaultRoutes(r, vaultHandlers, jwtAuthConfig)
	return r
}
