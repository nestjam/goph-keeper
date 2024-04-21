package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/nestjam/goph-keeper/internal/auth/repository/inmemory"
	"github.com/nestjam/goph-keeper/internal/auth/service"
)

func TestRegister(t *testing.T) {
	t.Run("regiser new user", func(t *testing.T) {
		const (
			email    = "user@email.com"
			password = "1234"
		)
		repo := inmemory.NewUserRepository()
		service := service.NewAuthService(repo)
		sut := NewAuthHandlers(service)
		r := newRegisterUserRequest(t, email, password)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		_, err := repo.FindByEmail(email)
		require.NoError(t, err)
	})
	t.Run("json is invalid", func(t *testing.T) {
		repo := inmemory.NewUserRepository()
		service := service.NewAuthService(repo)
		sut := NewAuthHandlers(service)
		r := newRegisterUserInvalidRequest(t)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("register failed", func(t *testing.T) {
		const (
			email    = "user@email.com"
			password = "1234"
		)
		service := &service.FakeAuthService{}
		service.RegisterFunc = func(user *model.User) (*model.User, error) {
			return nil, errors.New("failed to register")
		}
		sut := NewAuthHandlers(service)
		r := newRegisterUserRequest(t, email, password)
		w := httptest.NewRecorder()

		sut.Register().ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func newRegisterUserInvalidRequest(t *testing.T) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{{invalid}"))
}

func newRegisterUserRequest(t *testing.T, email, password string) *http.Request {
	t.Helper()

	r := RegisterUserRequest{
		Email:    email,
		Password: password,
	}
	body, err := json.Marshal(r)
	require.NoError(t, err)
	return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
}
