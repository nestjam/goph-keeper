package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	vaultHttp "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestGetSecretCommand(t *testing.T) {
	t.Run("get secret", func(t *testing.T) {
		const secretID = "11"
		wantSecret := vaultHttp.Secret{
			ID:   secretID,
			Data: "data",
		}
		wantURL := "/secrets/" + secretID
		wantCookie := &http.Cookie{
			Name: "jwt",
		}
		var gotURL string
		var gotCookie *http.Cookie
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = r.URL.String()
			gotCookie = findCookie(r.Cookies(), "jwt")
			resp := vaultHttp.GetSecretResponse{
				Secret: wantSecret,
			}
			_ = writeJSON(w, http.StatusOK, resp)
		}))
		defer server.Close()
		sut := newGetSecretCommand(secretID, server.URL, wantCookie)

		got := sut.execute()

		assert.Equal(t, wantURL, gotURL)
		assert.Equal(t, wantCookie, gotCookie)
		want := getSecretCompletedMsg{wantSecret}
		assert.Equal(t, want, got)
	})
	t.Run("invalid server address", func(t *testing.T) {
		sut := getSecretCommand{
			address: string([]byte{0x7f}), // ASCII control character
		}

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("failed to connect server", func(t *testing.T) {
		const secretID = "1"
		server := httptest.NewServer(nil)
		serverURL := server.URL
		server.Close()
		sut := newGetSecretCommand(secretID, serverURL, &http.Cookie{})

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("request is not successful", func(t *testing.T) {
		const secretID = "1"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		sut := newGetSecretCommand(secretID, server.URL, &http.Cookie{})

		got := sut.execute()

		msg, ok := got.(getSecretFailedMsg)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, msg.statusCode)
	})
}
