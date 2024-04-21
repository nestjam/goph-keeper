package http

import "github.com/go-chi/chi"

func MapAuthRoutes(r chi.Router, h *AuthHandlers) {
	r.Group(func(r chi.Router) {
		r.Post("/register", h.Register())
	})
}
