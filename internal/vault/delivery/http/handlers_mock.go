package http

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type vaultHandlersSpy struct {
	claims     map[string]interface{}
	callsCount int
}

func (m *vaultHandlersSpy) ListSecrets() http.HandlerFunc {
	return m.spy()
}

func (m *vaultHandlersSpy) AddSecret() http.HandlerFunc {
	return m.spy()
}

func (m *vaultHandlersSpy) spy() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.callsCount++
		_, m.claims, _ = jwtauth.FromContext(r.Context())
	})
}
