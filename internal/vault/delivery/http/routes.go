package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
)

func MapVaultRoutes(r chi.Router, h vault.VaultHandlers, cfg config.JWTAuthConfig) {
	const (
		secretsPath   = "/secrets"
		secretPattern = "/{secret}"
	)

	cookieBaker := utils.NewAuthCookieBaker(cfg)
	jwtAuth := cookieBaker.JWTAuth()

	r.Group(func(r chi.Router) {
		useJWTAuth(r, jwtAuth)

		r.Get(secretsPath, h.ListSecrets())
		r.Get(secretsPath+secretPattern, h.GetSecret())
		r.Delete(secretsPath+secretPattern, h.DeleteSecret())
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType(applicationJSON))
		useJWTAuth(r, jwtAuth)

		r.Post(secretsPath, h.AddSecret())
		r.Patch(secretsPath+secretPattern, h.UpdateSecret())
	})
}

func useJWTAuth(r chi.Router, jwtAuth *jwtauth.JWTAuth) {
	r.Use(jwtauth.Verifier(jwtAuth))
	r.Use(jwtauth.Authenticator(jwtAuth))
}
