package tui

import tea "github.com/charmbracelet/bubbletea"

type secretsModel struct {
}

func (m secretsModel) Init() tea.Cmd {
	return nil
}

func (m secretsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m secretsModel) View() string {
	return "secrets"
}
