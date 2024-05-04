package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nestjam/goph-keeper/internal/utils"
)

func TestLoginCommand(t *testing.T) {
	t.Run("login successful", func(t *testing.T) {
		gotURL := ""
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = r.URL.String()

			jwtCookie := &http.Cookie{
				Name: utils.JWTCookieName,
			}
			http.SetCookie(w, jwtCookie)

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		sut := loginCommand{
			address:  server.URL,
			email:    "user@email.com",
			password: "1234",
		}
		wantURL := "/login"

		got := sut.execute()

		assert.Equal(t, wantURL, gotURL)
		msg, ok := got.(loginCompletedMsg)
		assert.True(t, ok)
		assert.NotNil(t, msg.jwtCookie)
	})
	t.Run("invalid server address", func(t *testing.T) {
		sut := loginCommand{
			address:  string([]byte{0x7f}), // ASCII control character
			email:    "user@email.com",
			password: "1234",
		}

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("failed to connect server", func(t *testing.T) {
		server := httptest.NewServer(nil)
		serverURL := server.URL
		server.Close()
		sut := loginCommand{
			address:  serverURL,
			email:    "user@email.com",
			password: "1234",
		}

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("login is not successful", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		sut := loginCommand{
			address:  server.URL,
			email:    "user@email.com",
			password: "1234",
		}

		got := sut.execute()

		msg, ok := got.(loginFailedMsg)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, msg.statusCode)
	})
	t.Run("jwt cookie not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		sut := loginCommand{
			address:  server.URL,
			email:    "user@email.com",
			password: "1234",
		}

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
}
