package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
)

func MapVaultRoutes(r chi.Router, h vault.VaultHandlers, cfg config.JWTAuthConfig) {
	cookieBaker := utils.NewAuthCookieBaker(cfg)
	jwtAuth := cookieBaker.JWTAuth()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(jwtauth.Authenticator(jwtAuth))

		r.Get("/list", h.ListSecrets())
	})
}
