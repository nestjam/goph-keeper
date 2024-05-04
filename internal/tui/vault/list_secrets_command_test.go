package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	vaultHttp "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestListSecretsCommand(t *testing.T) {
	t.Run("list user secrets", func(t *testing.T) {
		want := []vaultHttp.Secret{
			{ID: "1"},
			{ID: "2"},
		}
		wantURL := "/secrets"
		wantCookie := &http.Cookie{
			Name: "jwt",
		}
		var gotURL string
		var gotCookie *http.Cookie
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = r.URL.String()
			gotCookie = findCookie(r.Cookies(), "jwt")
			resp := vaultHttp.ListSecretsResponse{
				List: want,
			}
			_ = writeJSON(w, http.StatusOK, resp)
		}))
		defer server.Close()
		sut := listSecretsCommand{
			address:   server.URL,
			jwtCookie: wantCookie,
		}

		msg := sut.Execute()

		assert.Equal(t, wantURL, gotURL)
		assert.Equal(t, wantCookie, gotCookie)
		res, ok := msg.(listSecretsCompletedMsg)
		assert.True(t, ok)
		got := res.secrets
		assert.Equal(t, want, got)
	})
	t.Run("invalid server address", func(t *testing.T) {
		sut := listSecretsCommand{
			address: string([]byte{0x7f}), // ASCII control character
		}

		got := sut.Execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("failed to connect server", func(t *testing.T) {
		server := httptest.NewServer(nil)
		serverURL := server.URL
		server.Close()
		sut := listSecretsCommand{
			address:   serverURL,
			jwtCookie: &http.Cookie{},
		}

		got := sut.Execute()

		assert.IsType(t, errMsg{}, got)
	})
	t.Run("request is not successful", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		sut := listSecretsCommand{
			address:   server.URL,
			jwtCookie: &http.Cookie{},
		}

		got := sut.Execute()

		msg, ok := got.(listSecretsFailedMsg)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, msg.statusCode)
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) error {
	const op = "write json"

	content, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, op)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(content)
	return nil
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == name {
			return cookies[i]
		}
	}

	return nil
}
