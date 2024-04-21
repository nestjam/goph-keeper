package server

import (
	"net/http"

	"github.com/go-chi/chi"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
)

func (s *Server) mapHandlers() http.Handler {
	r := chi.NewRouter()

	authHandlers := httpAuth.NewAuthHandlers()

	httpAuth.MapAuthRoutes(r, authHandlers)

	return r
}
