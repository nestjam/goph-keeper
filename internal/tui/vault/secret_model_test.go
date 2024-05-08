package vault

import (
	"errors"
	"net/http"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSecretModel_Init(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
	)

	sut := NewSecretModel(address, jwtCookie)

	got := sut.Init()

	assert.Nil(t, got)
}

func TestSecretModel_Update(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
	)

	t.Run("user exited by ctrl+c", func(t *testing.T) {
		sut := NewSecretModel(address, jwtCookie)
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("get secret request completed", func(t *testing.T) {
		sut := tea.Model(NewSecretModel(address, jwtCookie))
		want := httpVault.Secret{ID: "1", Data: "data"}
		msg := getSecretCompletedMsg{
			secret: want,
		}

		model, cmd := sut.Update(msg)

		m, _ := model.(secretModel)
		got := m.secret
		assert.Equal(t, want, got)
		assert.Nil(t, cmd)
	})
	t.Run("error on get secret", func(t *testing.T) {
		sut := NewSecretModel(address, jwtCookie)
		msg := errMsg{errors.New("error")}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		assert.Equal(t, msg.err, got.err)
	})
	t.Run("failed to get secret", func(t *testing.T) {
		sut := NewSecretModel(address, jwtCookie)
		const want = http.StatusBadRequest
		msg := getSecretFailedMsg{statusCode: want}

		model, _ := sut.Update(msg)

		m, _ := model.(secretModel)
		got := m.failtureStatusCode
		assert.Equal(t, want, got)
	})
	t.Run("return to list of secrets on esc", func(t *testing.T) {
		sut := NewSecretModel(address, jwtCookie)
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		model, cmd := sut.Update(msg)

		_, ok := model.(secretsModel)
		assert.True(t, ok)
		listSecretsCommand := NewListSecretsCommand(address, jwtCookie)
		assertEqualCmd(t, listSecretsCommand.Execute, cmd)
	})
	t.Run("create new secret requested", func(t *testing.T) {
		want := httpVault.Secret{}
		msg := createSecretRequestedMsg{}
		sut := tea.Model(NewSecretModel(address, jwtCookie))

		model, cmd := sut.Update(msg)

		m, _ := model.(secretModel)
		got := m.secret
		assert.Equal(t, want, got)
		assert.True(t, m.isNew)
		assert.Nil(t, cmd)
	})
	t.Run("save secret by ctrl+s", func(t *testing.T) {
		sut := NewSecretModel(address, jwtCookie)
		secret := httpVault.Secret{}
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
		want := httpVault.Secret{ID: "1", Data: "data"}
		msg := saveSecretCompletedMsg{want}
		sut := NewSecretModel(address, jwtCookie)
		sut.isNew = true

		model, cmd := sut.Update(msg)

		got, ok := model.(secretModel)
		assert.True(t, ok)
		assert.Nil(t, cmd)
		assert.Equal(t, want, got.secret)
		assert.Equal(t, want.Data, got.textarea.Value())
		assert.False(t, got.isNew)
	})
	t.Run("window size changed", func(t *testing.T) {
		sut := NewSecretModel(address, jwtCookie)
		msg := tea.WindowSizeMsg{Width: 100}
		require.NotEqual(t, msg.Width, sut.help.Width)

		model, _ := sut.Update(msg)

		got, ok := model.(secretModel)
		assert.True(t, ok)
		assert.Equal(t, msg.Width, got.help.Width)
	})
}
