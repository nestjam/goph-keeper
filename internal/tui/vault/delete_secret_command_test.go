package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteSecretCommand(t *testing.T) {
	t.Run("delete secret", func(t *testing.T) {
		const secretID = "1"
		wantURL := "/secrets/" + secretID
		wantCookie := &http.Cookie{
			Name: "jwt",
		}
		want := deleteSecretCompletedMsg{secretID}
		var gotURL string
		var gotCookie *http.Cookie
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = r.URL.String()
			gotCookie = findCookie(r.Cookies(), "jwt")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		sut := newDeleteSecretCommand(secretID, server.URL, wantCookie)

		got := sut.execute()

		assert.Equal(t, wantURL, gotURL)
		assert.Equal(t, wantCookie, gotCookie)
		assert.Equal(t, want, got)
	})
	t.Run("invalid server address", func(t *testing.T) {
		sut := deleteSecretCommand{
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
		sut := newDeleteSecretCommand(secretID, serverURL, &http.Cookie{})

		got := sut.execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("delete secret failed", func(t *testing.T) {
		const secretID = "1"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()
		sut := newDeleteSecretCommand(secretID, server.URL, &http.Cookie{})

		got := sut.execute()

		msg, ok := got.(deleteSecretFailedMsg)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, msg.statusCode)
	})
}
