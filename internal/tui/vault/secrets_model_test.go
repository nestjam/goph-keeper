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
		getSecretCommand := NewGetSecretCommand(wantID, address, jwtCookie)
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
}

func assertEqualCmd(t *testing.T, want, got tea.Cmd) {
	t.Helper()

	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)
	assert.Equal(t, gotValue.Pointer(), wantValue.Pointer())
}
