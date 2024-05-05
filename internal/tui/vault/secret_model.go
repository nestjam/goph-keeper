package vault

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type secretModel struct {
	textarea textarea.Model
	secret   httpVault.Secret
}

func NewSecretModel() secretModel {
	ti := textarea.New()
	ti.Focus()

	return secretModel{textarea: ti}
}

func (m secretModel) Init() tea.Cmd {
	return nil
}

func (m secretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
		}
	case getSecretCompletedMsg:
		{
			m.secret = msg.secret
			m.textarea.SetValue(m.secret.Data)
		}
	default:
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m secretModel) View() string {
	s := strings.Builder{}

	s.WriteString(m.secret.ID + "\n\n")

	s.WriteString(m.textarea.View())

	return s.String()
}
