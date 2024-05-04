package vault

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSecretsModel_Init(t *testing.T) {
	sut := NewSecretsModel()

	got := sut.Init()

	assert.Nil(t, got)
}

func TestSecretsModel_Update(t *testing.T) {
	t.Run("user exited by ctrl+c", func(t *testing.T) {
		sut := NewSecretsModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("user exited by esc", func(t *testing.T) {
		sut := NewSecretsModel()
		msg := tea.KeyMsg{Type: tea.KeyEsc}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("get secrets request completed", func(t *testing.T) {
		sut := tea.Model(NewSecretsModel())
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
}

func assertEqualCmd(t *testing.T, want, got tea.Cmd) {
	t.Helper()

	gotValue := reflect.ValueOf(got)
	wantValue := reflect.ValueOf(want)
	assert.Equal(t, gotValue.Pointer(), wantValue.Pointer())
}
