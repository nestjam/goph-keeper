package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	"github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/auth/service"
	"github.com/nestjam/goph-keeper/internal/config"
)

func (s *Server) mapHandlers() http.Handler {
	config := config.JWTAuthConfig{
		SignKey:       "supersecret",
		TokenExpiryIn: time.Hour,
	}
	authRepo := inmemory.NewUserRepository()
	authService := service.NewAuthService(authRepo)
	authHandlers := httpAuth.NewAuthHandlers(authService, config)

	r := chi.NewRouter()
	httpAuth.MapAuthRoutes(r, authHandlers)
	return r
}
