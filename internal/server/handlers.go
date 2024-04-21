package server

import (
	"net/http"

	"github.com/go-chi/chi"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	"github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/auth/service"
)

func (s *Server) mapHandlers() http.Handler {
	authRepo := inmemory.NewUserRepository()
	authService := service.NewAuthService(authRepo)
	authHandlers := httpAuth.NewAuthHandlers(authService)

	r := chi.NewRouter()
	httpAuth.MapAuthRoutes(r, authHandlers)
	return r
}
