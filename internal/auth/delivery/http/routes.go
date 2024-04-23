package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	applicationJSON   = "application/json"
	contentTypeHeader = "Content-Type"
)

func MapAuthRoutes(r chi.Router, h *AuthHandlers) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType(applicationJSON))

		r.Post("/register", h.Register())
	})
}
