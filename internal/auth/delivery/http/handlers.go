package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/nestjam/goph-keeper/internal/config"
	"github.com/nestjam/goph-keeper/internal/utils"
)

const (
	applicationJSON   = "application/json"
	contentTypeHeader = "Content-Type"
)

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandlers struct {
	service    auth.AuthService
	authConfig config.JWTAuthConfig
}

func NewAuthHandlers(service auth.AuthService, authConfig config.JWTAuthConfig) *AuthHandlers {
	return &AuthHandlers{
		service:    service,
		authConfig: authConfig,
	}
}

func (h *AuthHandlers) Register() http.HandlerFunc {
	cookieBaker := utils.NewAuthCookieBaker(h.authConfig)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUser(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		createdUser, err := h.service.Register(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie, err := cookieBaker.BakeCookie(createdUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusCreated)
	})
}

func getUser(r io.Reader) (*model.User, error) {
	const op = "get user"
	var userRequest RegisterUserRequest
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&userRequest)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	user := &model.User{
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}
	return user, nil
}
