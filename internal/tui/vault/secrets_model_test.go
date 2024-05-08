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

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSecretsModel_Init(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
	)
	sut := NewSecretsModel(address, jwtCookie)

	got := sut.Init()

	assert.Nil(t, got)
}

func TestSecretsModel_Update(t *testing.T) {
	var (
		address   = "/"
		jwtCookie = &http.Cookie{}
	)

	t.Run("user exited by ctrl+c", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("user exited by esc", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("get secrets request completed", func(t *testing.T) {
		sut := tea.Model(NewSecretsModel(address, jwtCookie))
		want := []httpVault.Secret{
			{ID: "2"},
			{ID: "3"},
		}
		wantRows := []table.Row{
			{"1", "2"},
			{"2", "3"},
		}
		msg := listSecretsCompletedMsg{
			secrets: want,
		}

		model, cmd := sut.Update(msg)

		m, _ := model.(secretsModel)
		got := m.secrets
		assert.Equal(t, want, got)
		assert.Nil(t, cmd)
		assert.Equal(t, wantRows, m.table.Rows())
	})
	t.Run("user pressed enter on selected row", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		secrets := []httpVault.Secret{
			{ID: "2"},
			{ID: "3"},
		}
		sut.secrets = secrets
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
	t.Run("error on get secret", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		msg := errMsg{errors.New("error")}

		model, _ := sut.Update(msg)

		got, _ := model.(secretsModel)
		assert.Equal(t, msg.err, got.err)
	})
	t.Run("failed to list secrets", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		const want = http.StatusBadRequest
		msg := listSecretsFailedMsg{statusCode: want}

		model, _ := sut.Update(msg)

		m, _ := model.(secretsModel)
		got := m.failtureStatusCode
		assert.Equal(t, want, got)
	})
	t.Run("delete selected secret on del", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		secrets := []httpVault.Secret{
			{ID: "2"},
		}
		sut.secrets = secrets
		rows := []table.Row{
			{"1", "2"},
		}
		sut.table.SetRows(rows)
		sut.table.GotoTop()
		const wantID = "2"
		msg := tea.KeyMsg{Type: tea.KeyDelete}

		model, cmd := sut.Update(msg)

		_, ok := model.(secretsModel)
		assert.True(t, ok)
		deleteSecretCommand := newDeleteSecretCommand(wantID, address, jwtCookie)
		assertEqualCmd(t, deleteSecretCommand.execute, cmd)
	})
	t.Run("selected secret deleted", func(t *testing.T) {
		const secretID = "2"
		secrets := []httpVault.Secret{
			{ID: secretID},
		}
		sut := NewSecretsModel(address, jwtCookie)
		sut.secrets = secrets
		rows := []table.Row{
			{"1", secretID},
		}
		sut.table.SetRows(rows)
		sut.table.GotoTop()
		msg := deleteSecretCompletedMsg{secretID}

		model, _ := sut.Update(msg)

		m, ok := model.(secretsModel)
		assert.True(t, ok)
		assert.Empty(t, len(m.table.Rows()))
	})
	t.Run("add new secret by ctrl+n", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		msg := tea.KeyMsg{Type: tea.KeyCtrlN}

		model, cmd := sut.Update(msg)

		_, ok := model.(secretModel)
		assert.True(t, ok)
		createSecretCommand := newCreateSecretCommand()
		assertEqualCmd(t, createSecretCommand.execute, cmd)
	})
	t.Run("window size changed", func(t *testing.T) {
		sut := NewSecretsModel(address, jwtCookie)
		msg := tea.WindowSizeMsg{Width: 100}
		require.NotEqual(t, msg.Width, sut.help.Width)

		model, _ := sut.Update(msg)

		got, ok := model.(secretsModel)
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
