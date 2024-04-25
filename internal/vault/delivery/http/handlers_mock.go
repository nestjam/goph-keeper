package http

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type vaultHandlersMock struct {
	claims map[string]interface{}
	calls  int
}

func (m *vaultHandlersMock) ListSecrets() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.calls++
		_, m.claims, _ = jwtauth.FromContext(r.Context())
	})
}
