package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/config"
)

func TestList(t *testing.T) {
	config := config.JWTAuthConfig{
		SignKey:       "secret",
		TokenExpiryIn: time.Minute,
	}

	t.Run("empty list", func(t *testing.T) {
		sut := NewVaultHandlers(config)
		r := newListSecretsRequest(t, "/")
		w := httptest.NewRecorder()

		sut.ListSecrets().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assertContentType(t, applicationJSON, w)
		resp := getListDataResponse(t, w.Body)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.List)
	})
}

func newListSecretsRequest(t *testing.T, path string) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodGet, path, nil)
}

func getListDataResponse(t *testing.T, r io.Reader) *ListDataResponse {
	t.Helper()

	var resp ListDataResponse
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&resp)
	require.NoError(t, err)
	return &resp
}

func assertContentType(t *testing.T, want string, r *httptest.ResponseRecorder) {
	t.Helper()
	assert.Equal(t, want, r.Header().Get(contentTypeHeader))
}
