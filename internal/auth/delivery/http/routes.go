package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MapAuthRoutes(r chi.Router, h *AuthHandlers) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType(applicationJSON))

		r.Post("/register", h.Register())
		r.Post("/login", h.Login())
	})
}
