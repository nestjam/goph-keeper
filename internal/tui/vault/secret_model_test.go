package vault

import (
	"errors"
	"net/http"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

func TestSecretModel_Init(t *testing.T) {
	sut := NewSecretModel()

	got := sut.Init()

	assert.Nil(t, got)
}

func TestSecretModel_Update(t *testing.T) {
	t.Run("user exited by ctrl+c", func(t *testing.T) {
		sut := NewSecretModel()
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		_, cmd := sut.Update(msg)

		assertEqualCmd(t, tea.Quit, cmd)
	})
	t.Run("get secret request completed", func(t *testing.T) {
		sut := tea.Model(NewSecretModel())
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
		sut := NewSecretModel()
		msg := errMsg{errors.New("error")}

		model, _ := sut.Update(msg)

		got, _ := model.(secretModel)
		assert.Equal(t, msg.err, got.err)
	})
	t.Run("failed to get secret", func(t *testing.T) {
		sut := NewSecretModel()
		const want = http.StatusBadRequest
		msg := getSecretFailedMsg{statusCode: want}

		model, _ := sut.Update(msg)

		m, _ := model.(secretModel)
		got := m.failtureStatusCode
		assert.Equal(t, want, got)
	})
}
