package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/pkg/errors"
)

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandlers struct {
	service auth.AuthService
}

func NewAuthHandlers(service auth.AuthService) *AuthHandlers {
	return &AuthHandlers{
		service: service,
	}
}

func (h *AuthHandlers) Register() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUser(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = h.service.Register(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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
