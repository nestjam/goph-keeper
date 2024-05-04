package vault

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type secretsModel struct {
	secrets []httpVault.Secret

	table table.Model
}

func NewSecretsModel() secretsModel {
	const (
		numWidth    = 4
		idWidth     = 60
		tableHeight = 7
	)
	columns := []table.Column{
		{Title: "#", Width: numWidth},
		{Title: "ID", Width: idWidth},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return secretsModel{
		table: t,
	}
}

func (m secretsModel) Init() tea.Cmd {
	return nil
}

func (m secretsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	case listSecretsCompletedMsg:
		{
			m.secrets = msg.secrets
			rows := make([]table.Row, len(m.secrets))
			for i := 0; i < len(m.secrets); i++ {
				secret := m.secrets[i]
				rows[i] = table.Row{strconv.Itoa(i + 1), secret.ID}
			}
			m.table.SetRows(rows)
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m secretsModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
