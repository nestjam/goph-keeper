package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type secretsModel struct {
	secrets []httpVault.Secret
}

func (m secretsModel) Init() tea.Cmd {
	return nil
}

func (m secretsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case listSecretsCompletedMsg:
		{
			m.secrets = msg.secrets
		}
	default:
	}

	return m, nil
}

func (m secretsModel) View() string {
	s := strings.Builder{}

	s.WriteString("secrets: ")

	if len(m.secrets) > 0 {
		s.WriteString(strconv.Itoa(len(m.secrets)))
	}

	return s.String()
}
