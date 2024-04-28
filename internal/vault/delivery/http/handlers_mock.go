package http

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type vaultHandlersSpy struct {
	claims                 map[string]interface{}
	listSecretsCallsCount  int
	addSecretCallsCount    int
	getSecretCallsCount    int
	deleteSecretCallsCount int
}

func (m *vaultHandlersSpy) ListSecrets() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.listSecretsCallsCount++
		_, m.claims, _ = jwtauth.FromContext(r.Context())
	})
}

func (m *vaultHandlersSpy) AddSecret() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.addSecretCallsCount++
		_, m.claims, _ = jwtauth.FromContext(r.Context())
	})
}

func (m *vaultHandlersSpy) GetSecret() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.getSecretCallsCount++
	})
}

func (m *vaultHandlersSpy) DeleteSecret() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.deleteSecretCallsCount++
	})
}
