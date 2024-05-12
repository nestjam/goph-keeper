package vault

import (
	"errors"
	"net/http"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/tui/vault/cache"
	vault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSecretModel_Init(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
	)
	cache := cache.New()
	sut := NewSecretModel(address, jwtCookie, cache)

	got := sut.Init()

	assert.Nil(t, got)
}

func TestSecretModel_Update(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
	)

	t.Run("user enter text", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		msg := tea.KeyMsg{Type: tea.KeySpace, Runes: []rune("text")}

		model, cmd := sut.Update(msg)

		_, ok := model.(secretModel)
		assert.True(t, ok)
		assert.NotNil(t, cmd)
	})
	t.Run("user exited by ctrl+c", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("get secret request completed", func(t *testing.T) {
		cache := cache.New()
		sut := tea.Model(NewSecretModel(address, jwtCookie, cache))
		wantSecret := vault.Secret{ID: "1", Data: "data"}
		msg := getSecretCompletedMsg{
			secret: wantSecret,
		}

		model, cmd := sut.Update(msg)

		got, _ := model.(secretModel)
		gotSecret := got.secret
		assert.Equal(t, wantSecret, gotSecret)
		assert.Nil(t, cmd)
		cachedSecret, dataCached, ok := cache.GetSecret(wantSecret.ID)
		require.True(t, ok)
		require.True(t, dataCached)
		assert.Equal(t, wantSecret, *cachedSecret)
	})
	t.Run("error on get secret", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		msg := getSecretFailedMsg{err: errors.New("error")}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		assert.Equal(t, msg.err, got.err)
		assert.True(t, got.isOffline)
		assert.False(t, got.keys.Save.Enabled())
	})
	t.Run("failed to get secret", func(t *testing.T) {
		secret := &vault.Secret{
			ID:   "1",
			Data: "data",
		}
		cache := cache.New()
		cache.CacheSecret(secret)
		sut := NewSecretModel(address, jwtCookie, cache)
		const want = http.StatusBadRequest
		msg := getSecretFailedMsg{
			statusCode: want,
			secretID:   secret.ID,
		}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		gotStatusCode := got.failtureStatusCode
		assert.Equal(t, want, gotStatusCode)
		assert.True(t, got.isOffline)
		assert.False(t, got.keys.Save.Enabled())
		assert.Equal(t, secret.Data, got.textarea.Value())
		assert.True(t, got.dataCached)
	})
	t.Run("clear text if failed to get secret and secret data is not cached", func(t *testing.T) {
		secrets := []*vault.Secret{
			{ID: "1"},
		}
		cache := cache.New()
		cache.CacheSecrets(secrets)
		sut := NewSecretModel(address, jwtCookie, cache)
		sut.textarea.SetValue("text")
		const want = http.StatusBadRequest
		msg := getSecretFailedMsg{
			secretID:   secrets[0].ID,
			statusCode: want,
		}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		assert.Empty(t, got.textarea.Value())
		assert.False(t, got.dataCached)
		assert.Equal(t, noCachedData, got.textarea.Placeholder)
	})
	t.Run("clear text if failed to get secret and secret is not cached", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		sut.textarea.SetValue("text")
		const want = http.StatusBadRequest
		msg := getSecretFailedMsg{statusCode: want}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		assert.Empty(t, got.textarea.Value())
		assert.False(t, got.dataCached)
	})
	t.Run("error", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		msg := errMsg{err: errors.New("error")}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		assert.Equal(t, msg.err, got.err)
		assert.True(t, got.isOffline)
		assert.False(t, got.keys.Save.Enabled())
	})
	t.Run("view cached secret when failed to get secret", func(t *testing.T) {
		cache := cache.New()
		secret := &vault.Secret{ID: "1", Data: "123"}
		cache.CacheSecret(secret)
		sut := NewSecretModel(address, jwtCookie, cache)
		const want = http.StatusBadRequest
		msg := getSecretFailedMsg{
			statusCode: want,
			secretID:   secret.ID,
		}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		assert.Equal(t, secret.Data, got.textarea.Value())
		assert.Equal(t, *secret, got.secret)
	})
	t.Run("return to list of secrets on esc", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		model, cmd := sut.Update(msg)

		_, ok := model.(SecretsModel)
		assert.True(t, ok)
		listSecretsCommand := NewListSecretsCommand(address, jwtCookie)
		assertEqualCmd(t, listSecretsCommand.Execute, cmd)
	})
	t.Run("create new secret requested", func(t *testing.T) {
		want := vault.Secret{}
		msg := createSecretRequestedMsg{}
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)

		model, cmd := sut.Update(msg)

		m, _ := model.(secretModel)
		got := m.secret
		assert.Equal(t, want, got)
		assert.True(t, m.isNew)
		assert.Nil(t, cmd)
	})
	t.Run("save secret by ctrl+s", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		secret := vault.Secret{}
		sut.secret = secret
		sut.textarea.SetValue("data")
		msg := tea.KeyMsg{Type: tea.KeyCtrlS}

		model, cmd := sut.Update(msg)

		_, ok := model.(secretModel)
		assert.True(t, ok)
		saveSecretCommand := newSaveSecretCommand(secret, address, jwtCookie)
		assertEqualCmd(t, saveSecretCommand.execute, cmd)
	})
	t.Run("save secret completed", func(t *testing.T) {
		want := vault.Secret{ID: "1", Data: "data"}
		msg := saveSecretCompletedMsg{want}
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		sut.isNew = true

		model, cmd := sut.Update(msg)

		got, ok := model.(secretModel)
		assert.True(t, ok)
		assert.Nil(t, cmd)
		assert.Equal(t, want, got.secret)
		assert.Equal(t, want.Data, got.textarea.Value())
		assert.False(t, got.isNew)
		cachedSecret, _, ok := cache.GetSecret(want.ID)
		assert.True(t, ok)
		assert.Equal(t, want, *cachedSecret)
	})
	t.Run("window size changed", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretModel(address, jwtCookie, cache)
		msg := tea.WindowSizeMsg{Width: 100}
		require.NotEqual(t, msg.Width, sut.help.Width)

		model, _ := sut.Update(msg)

		got, ok := model.(secretModel)
		assert.True(t, ok)
		assert.Equal(t, msg.Width, got.help.Width)
	})
}
