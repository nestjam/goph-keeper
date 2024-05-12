package vault

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/tui/vault/cache"
	vault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSecretsModel_Init(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
		cache     = cache.New()
	)
	sut := NewSecretsModel(address, jwtCookie, cache)

	got := sut.Init()

	assert.Nil(t, got)
}

func TestSecretsModel_Update(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
	)

	t.Run("user exited by ctrl+c", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretsModel(address, jwtCookie, cache)
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("user exited by esc", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretsModel(address, jwtCookie, cache)
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("get secrets request completed", func(t *testing.T) {
		cache := cache.New()
		sut := tea.Model(NewSecretsModel(address, jwtCookie, cache))
		wantSecrets := []*vault.Secret{
			{ID: "2"},
			{ID: "3"},
		}
		wantRows := []table.Row{
			{"1", "2"},
			{"2", "3"},
		}
		msg := listSecretsCompletedMsg{
			secrets: wantSecrets,
		}

		model, cmd := sut.Update(msg)

		got, _ := model.(SecretsModel)
		gotSecrets := got.cache.ListSecrets()
		assert.ElementsMatch(t, wantSecrets, gotSecrets)
		assert.Nil(t, cmd)
		assert.ElementsMatch(t, wantRows, got.table.Rows())
		assert.False(t, got.isOffline)
	})
	t.Run("user pressed enter on selected row", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretsModel(address, jwtCookie, cache)
		secrets := []*vault.Secret{
			{ID: "2"},
			{ID: "3"},
		}
		cache.CacheSecrets(secrets)
		rows := []table.Row{
			{"1", "2"},
			{"2", "3"},
		}
		sut.table.SetRows(rows)
		sut.table.GotoTop()
		const wantID = "2"
		require.Equal(t, wantID, sut.table.SelectedRow()[1])
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		model, cmd := sut.Update(msg)

		_, ok := model.(secretModel)
		assert.True(t, ok)
		getSecretCommand := newGetSecretCommand(wantID, address, jwtCookie)
		assert.Equal(t, wantID, getSecretCommand.secretID)
		assertEqualCmd(t, getSecretCommand.execute, cmd)
	})
	t.Run("error on get secrets", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretsModel(address, jwtCookie, cache)
		msg := listSecretsFailedMsg{err: errors.New("error")}

		model, _ := sut.Update(msg)

		got, _ := model.(SecretsModel)
		assert.Equal(t, msg.err, got.err)
		assert.Equal(t, zeroStatusCode, got.failtureStatusCode)
		assert.True(t, got.isOffline)
		assert.False(t, got.keys.Add.Enabled())
		assert.False(t, got.keys.Delete.Enabled())
	})
	t.Run("failed to list secrets", func(t *testing.T) {
		wantRows := []table.Row{
			{"1", "2"},
		}
		secrets := []*vault.Secret{
			{ID: "2"},
		}
		cache := cache.New()
		cache.CacheSecrets(secrets)
		sut := NewSecretsModel(address, jwtCookie, cache)
		const wantStatusCode = http.StatusBadRequest
		msg := listSecretsFailedMsg{statusCode: wantStatusCode}

		model, _ := sut.Update(msg)

		got, _ := model.(SecretsModel)
		assert.Nil(t, got.err)
		assert.Equal(t, wantStatusCode, got.failtureStatusCode)
		assert.True(t, got.isOffline)
		assert.False(t, got.keys.Add.Enabled())
		assert.False(t, got.keys.Delete.Enabled())
		assert.ElementsMatch(t, wantRows, got.table.Rows())
	})
	t.Run("success retry after failed to list secrets", func(t *testing.T) {
		cache := cache.New()
		sut := tea.Model(NewSecretsModel(address, jwtCookie, cache))
		secrets := []*vault.Secret{}
		var msg tea.Msg = listSecretsFailedMsg{statusCode: http.StatusTooManyRequests}
		sut, _ = sut.Update(msg)

		msg = listSecretsCompletedMsg{
			secrets: secrets,
		}
		model, _ := sut.Update(msg)

		got, _ := model.(SecretsModel)
		assert.False(t, got.isOffline)
		assert.Equal(t, zeroStatusCode, got.failtureStatusCode)
		assert.True(t, got.keys.Add.Enabled())
		assert.True(t, got.keys.Delete.Enabled())
	})
	t.Run("delete selected secret on del", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretsModel(address, jwtCookie, cache)
		secrets := []*vault.Secret{
			{ID: "2"},
		}
		cache.CacheSecrets(secrets)
		rows := []table.Row{
			{"1", "2"},
		}
		sut.table.SetRows(rows)
		sut.table.GotoTop()
		const wantID = "2"
		msg := tea.KeyMsg{Type: tea.KeyDelete}

		model, cmd := sut.Update(msg)

		_, ok := model.(SecretsModel)
		assert.True(t, ok)
		deleteSecretCommand := newDeleteSecretCommand(wantID, address, jwtCookie)
		assertEqualCmd(t, deleteSecretCommand.execute, cmd)
	})
	t.Run("selected secret deleted", func(t *testing.T) {
		const secretID = "2"
		secrets := []*vault.Secret{
			{ID: secretID},
		}
		cache := cache.New()
		cache.CacheSecrets(secrets)
		sut := NewSecretsModel(address, jwtCookie, cache)
		rows := []table.Row{
			{"1", secretID},
		}
		sut.table.SetRows(rows)
		sut.table.GotoTop()
		msg := deleteSecretCompletedMsg{secretID}

		model, _ := sut.Update(msg)

		m, ok := model.(SecretsModel)
		assert.True(t, ok)
		assert.Empty(t, len(m.table.Rows()))
		cachedSecrets := cache.ListSecrets()
		assert.Empty(t, cachedSecrets)
	})
	t.Run("add new secret by ctrl+n", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretsModel(address, jwtCookie, cache)
		msg := tea.KeyMsg{Type: tea.KeyCtrlN}

		model, cmd := sut.Update(msg)

		_, ok := model.(secretModel)
		assert.True(t, ok)
		createSecretCommand := newCreateSecretCommand()
		assertEqualCmd(t, createSecretCommand.execute, cmd)
	})
	t.Run("window size changed", func(t *testing.T) {
		cache := cache.New()
		sut := NewSecretsModel(address, jwtCookie, cache)
		msg := tea.WindowSizeMsg{Width: 100}
		require.NotEqual(t, msg.Width, sut.help.Width)

		model, _ := sut.Update(msg)

		got, ok := model.(SecretsModel)
		assert.True(t, ok)
		assert.Equal(t, msg.Width, got.help.Width)
	})
}

func assertEqualCmd(t *testing.T, want, got tea.Cmd) {
	t.Helper()

	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)
	assert.Equal(t, gotValue.Pointer(), wantValue.Pointer())
}
