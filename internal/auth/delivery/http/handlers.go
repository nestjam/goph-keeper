package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
)

const (
	jwtCookieName = "jwt"
	userIDClaim   = "user_id"
	jwtAlg        = "HS256"
)

type JWTAuthConfig struct {
	SignKey       string
	TokenExpiryIn time.Duration
}

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandlers struct {
	service    auth.AuthService
	authConfig JWTAuthConfig
}

func NewAuthHandlers(service auth.AuthService, authConfig JWTAuthConfig) *AuthHandlers {
	return &AuthHandlers{
		service:    service,
		authConfig: authConfig,
	}
}

func (h *AuthHandlers) Register() http.HandlerFunc {
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

		err = h.setAuthCookie(w, createdUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func (h *AuthHandlers) setAuthCookie(w http.ResponseWriter, user *model.User) error {
	const op = "set auth cookie"

	jwtAuth := jwtauth.New("HS256", []byte(h.authConfig.SignKey), nil)
	claims := make(map[string]interface{})
	claims[userIDClaim] = user.ID.String()
	jwtauth.SetExpiryIn(claims, h.authConfig.TokenExpiryIn)

	_, token, err := jwtAuth.Encode(claims)
	if err != nil {
		return errors.Wrap(err, op)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		MaxAge:   int(h.authConfig.TokenExpiryIn / time.Second),
		HttpOnly: true,
	})
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
