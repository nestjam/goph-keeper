package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/nestjam/goph-keeper/internal/utils"
)

func TestRegisterCommand(t *testing.T) {
	t.Run("register successful", func(t *testing.T) {
		gotURL := ""
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = r.URL.String()

			jwtCookie := &http.Cookie{
				Name: utils.JWTCookieName,
			}
			http.SetCookie(w, jwtCookie)

			w.WriteHeader(http.StatusCreated)
		}))
		defer server.Close()
		wantURL := "/register"
		client := resty.New()
		sut := newRegisterCommand(server.URL, "user@email.com", "1234", client)

		got := sut.execute()

		assert.Equal(t, wantURL, gotURL)
		msg, ok := got.(registerCompletedMsg)
		assert.True(t, ok)
		assert.NotNil(t, msg.jwtCookie)
	})
	t.Run("invalid server address", func(t *testing.T) {
		client := resty.New()
		sut := newRegisterCommand(string([]byte{0x7f}), // ASCII control character
			"user@email.com", "1234", client)

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("failed to connect server", func(t *testing.T) {
		server := httptest.NewServer(nil)
		serverURL := server.URL
		server.Close()
		client := resty.New()
		sut := newRegisterCommand(serverURL, "user@email.com", "1234", client)

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("register is not successful", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		client := resty.New()
		sut := newRegisterCommand(server.URL, "user@email.com", "1234", client)

		got := sut.execute()

		msg, ok := got.(registerFailedMsg)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, msg.statusCode)
	})
	t.Run("jwt cookie not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		client := resty.New()
		sut := newRegisterCommand(server.URL, "user@email.com", "1234", client)

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
}
