package http

import "net/http"

type AuthHandlers struct {
}

func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{}
}

func (h *AuthHandlers) Register() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
}
