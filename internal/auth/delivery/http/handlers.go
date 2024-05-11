package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
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

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandlers struct {
	service     auth.AuthService
	cookieBaker *utils.AuthCookieBaker
}

func NewAuthHandlers(service auth.AuthService, authConfig config.JWTAuthConfig) *AuthHandlers {
	return &AuthHandlers{
		service:     service,
		cookieBaker: utils.NewAuthCookieBaker(authConfig),
	}
}

//nolint:dupl //register method
func (h *AuthHandlers) Register() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUser(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		userID, err := h.service.Register(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = setAuthCookie(w, userID, h.cookieBaker)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

//nolint:dupl //login method
func (h *AuthHandlers) Login() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUser(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		userID, err := h.service.Login(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = setAuthCookie(w, userID, h.cookieBaker)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func setAuthCookie(w http.ResponseWriter, userID uuid.UUID, baker *utils.AuthCookieBaker) error {
	const op = "set auth cookie"

	cookie, err := baker.BakeCookie(userID)
	if err != nil {
		return errors.Wrap(err, op)
	}
	http.SetCookie(w, cookie)

	return nil
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
