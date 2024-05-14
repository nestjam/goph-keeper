package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	vaultHttp "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestListSecretsCommand(t *testing.T) {
	t.Run("list user secrets", func(t *testing.T) {
		wantSecrets := []*vaultHttp.Secret{
			{ID: "1"},
			{ID: "2"},
		}
		wantURL := "/secrets"
		wantCookie := &http.Cookie{
			Name: "auth",
		}
		var gotURL string
		var gotCookie *http.Cookie
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = r.URL.String()
			gotCookie = findCookie(r.Cookies(), "auth")
			resp := vaultHttp.ListSecretsResponse{
				List: []vaultHttp.Secret{*wantSecrets[0], *wantSecrets[1]},
			}
			_ = writeJSON(w, http.StatusOK, resp)
		}))
		defer server.Close()
		client := resty.New()
		sut := NewListSecretsCommand(server.URL, wantCookie, client)

		got := sut.Execute()

		assert.Equal(t, wantURL, gotURL)
		assert.Equal(t, wantCookie, gotCookie)
		want := listSecretsCompletedMsg{wantSecrets}
		assert.Equal(t, want, got)
	})
	t.Run("invalid server address", func(t *testing.T) {
		sut := listSecretsCommand{
			address: string([]byte{0x7f}), // ASCII control character
		}

		msg := sut.Execute()

		got, ok := msg.(listSecretsFailedMsg)
		assert.True(t, ok)
		assert.NotNil(t, got.err)
		assert.Equal(t, zeroStatusCode, got.statusCode)
	})
	t.Run("failed to connect server", func(t *testing.T) {
		server := httptest.NewServer(nil)
		serverURL := server.URL
		server.Close()
		client := resty.New()
		sut := NewListSecretsCommand(serverURL, &http.Cookie{}, client)

		msg := sut.Execute()

		got, ok := msg.(listSecretsFailedMsg)
		assert.True(t, ok)
		assert.NotNil(t, got.err)
		assert.Equal(t, zeroStatusCode, got.statusCode)
	})
	t.Run("request is not successful", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()
		client := resty.New()
		sut := NewListSecretsCommand(server.URL, &http.Cookie{}, client)

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
