package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSaveSecretCommand(t *testing.T) {
	t.Run("save new secret", func(t *testing.T) {
		const (
			wantMethod = "POST"
			secretID   = "1"
		)
		wantURL := "/secrets"
		wantCookie := &http.Cookie{
			Name: "jwt",
		}
		secret := httpVault.Secret{
			Data: "data",
		}
		want := saveSecretCompletedMsg{
			secret: httpVault.Secret{ID: secretID, Data: secret.Data},
		}
		var gotURL string
		var gotCookie *http.Cookie
		var gotMethod string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotURL = r.URL.String()
			gotCookie = findCookie(r.Cookies(), "jwt")
			s := httpVault.Secret{ID: secretID} // server does not return secret data
			_ = writeJSON(w, http.StatusCreated, httpVault.AddSecretResponse{Secret: s})
		}))
		defer server.Close()
		sut := newSaveSecretCommand(secret, server.URL, wantCookie)

		got := sut.execute()

		assert.Equal(t, wantURL, gotURL)
		assert.Equal(t, wantMethod, gotMethod)
		assert.Equal(t, wantCookie, gotCookie)
		assert.Equal(t, want, got)
	})
	t.Run("add secret: invalid server address", func(t *testing.T) {
		sut := saveSecretCommand{
			address: string([]byte{0x7f}), // ASCII control character
		}

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("add secret: failed to connect server", func(t *testing.T) {
		server := httptest.NewServer(nil)
		serverURL := server.URL
		server.Close()
		secret := httpVault.Secret{
			Data: "data",
		}
		sut := newSaveSecretCommand(secret, serverURL, &http.Cookie{})

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("add secret failed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()
		jwtCookie := &http.Cookie{}
		secret := httpVault.Secret{
			Data: "data",
		}
		sut := newSaveSecretCommand(secret, server.URL, jwtCookie)

		got := sut.execute()

		msg, ok := got.(saveSecretFailedMsg)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, msg.statusCode)
	})
	t.Run("save edited secret", func(t *testing.T) {
		const (
			secretID   = "1"
			wantMethod = "PATCH"
		)
		wantURL := "/secrets/" + secretID
		wantCookie := &http.Cookie{
			Name: "jwt",
		}
		secret := httpVault.Secret{
			ID:   secretID,
			Data: "data",
		}
		want := saveSecretCompletedMsg{
			secret: httpVault.Secret{
				ID:   secretID,
				Data: secret.Data,
			},
		}
		var gotURL string
		var gotCookie *http.Cookie
		var gotMethod string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotURL = r.URL.String()
			gotCookie = findCookie(r.Cookies(), "jwt")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		sut := newSaveSecretCommand(secret, server.URL, wantCookie)

		got := sut.execute()

		assert.Equal(t, wantURL, gotURL)
		assert.Equal(t, wantMethod, gotMethod)
		assert.Equal(t, wantCookie, gotCookie)
		assert.Equal(t, want, got)
	})
	t.Run("update secret: invalid server address", func(t *testing.T) {
		sut := saveSecretCommand{
			address: string([]byte{0x7f}), // ASCII control character
			secret:  httpVault.Secret{ID: "1"},
		}

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("update secret: failed to connect server", func(t *testing.T) {
		server := httptest.NewServer(nil)
		serverURL := server.URL
		server.Close()
		secret := httpVault.Secret{
			ID:   "1",
			Data: "data",
		}
		sut := newSaveSecretCommand(secret, serverURL, &http.Cookie{})

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("update secret failed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()
		jwtCookie := &http.Cookie{}
		secret := httpVault.Secret{
			ID:   "1",
			Data: "data",
		}
		sut := newSaveSecretCommand(secret, server.URL, jwtCookie)

		got := sut.execute()

		msg, ok := got.(saveSecretFailedMsg)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, msg.statusCode)
	})
}
