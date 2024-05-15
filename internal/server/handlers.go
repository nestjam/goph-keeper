package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	users "github.com/nestjam/goph-keeper/internal/auth/repository/pgsql"
	serviceAuth "github.com/nestjam/goph-keeper/internal/auth/service"
	"github.com/nestjam/goph-keeper/internal/config"
	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
	keys "github.com/nestjam/goph-keeper/internal/vault/repository/pgsql/key"
	secrets "github.com/nestjam/goph-keeper/internal/vault/repository/pgsql/secret"
	serviceVault "github.com/nestjam/goph-keeper/internal/vault/service"
)

func (s *Server) mapHandlers(ctx context.Context) (http.Handler, error) {
	const op = "map handlers"
	jwtAuthConfig := config.JWTAuthConfig{
		SignKey:       "supersecret",
		TokenExpiryIn: time.Hour,
	}

	authRepo, err := users.NewUserRepository(ctx, s.conf.Postgres.DataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	authService := serviceAuth.NewAuthService(authRepo)
	authHandlers := httpAuth.NewAuthHandlers(authService, jwtAuthConfig)

	secretRepo, err := secrets.NewSecretRepository(ctx, s.conf.Postgres.DataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	keyRepo, err := keys.NewDataKeyRepository(ctx, s.conf.Postgres.DataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	vaultService := serviceVault.NewVaultService(secretRepo, keyRepo, s.rootKey)
	vaultHandlers := httpVault.NewVaultHandlers(vaultService, jwtAuthConfig)

	r := chi.NewRouter()
	httpAuth.MapAuthRoutes(r, authHandlers)
	httpVault.MapVaultRoutes(r, vaultHandlers, jwtAuthConfig)
	return r, nil
}
