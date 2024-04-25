package http

import "net/http"

type vaultHandlersMock struct {
	listSecretsCalls int
}

func (m *vaultHandlersMock) ListSecrets() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.listSecretsCalls++
		w.WriteHeader(http.StatusNotImplemented)
	})
}
