package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSaveSecretCommand(t *testing.T) {
	t.Run("save secret", func(t *testing.T) {
		wantURL := "/secrets"
		wantCookie := &http.Cookie{
			Name: "jwt",
		}
		secret := httpVault.Secret{
			Data: "data",
		}
		want := saveSecretCompletedMsg{secret}
		var gotURL string
		var gotCookie *http.Cookie
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = r.URL.String()
			gotCookie = findCookie(r.Cookies(), "jwt")
			s := secretFromRequest(t, r)
			_ = writeJSON(w, http.StatusCreated, httpVault.AddSecretResponse{Secret: s})
		}))
		defer server.Close()
		sut := newSaveSecretCommand(secret, server.URL, wantCookie)

		got := sut.execute()

		assert.Equal(t, wantURL, gotURL)
		assert.Equal(t, wantCookie, gotCookie)
		assert.Equal(t, want, got)
	})
	t.Run("invalid server address", func(t *testing.T) {
		sut := saveSecretCommand{
			address: string([]byte{0x7f}), // ASCII control character
		}

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("failed to connect server", func(t *testing.T) {
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
}

func secretFromRequest(t *testing.T, r *http.Request) httpVault.Secret {
	t.Helper()

	var req httpVault.AddSecretRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	require.NoError(t, err)

	return req.Secret
}
