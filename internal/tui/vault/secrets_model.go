package vault

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type secretsModel struct {
	err                error
	jwtCookie          *http.Cookie
	address            string
	secrets            []httpVault.Secret
	table              table.Model
	failtureStatusCode int
}

func NewSecretsModel(address string, jwtCookie *http.Cookie) secretsModel {
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
		address:   address,
		jwtCookie: jwtCookie,
		table:     t,
	}
}

func (m secretsModel) Init() tea.Cmd {
	return nil
}

func (m secretsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			{
				const idColumnIndex = 1
				id := m.table.SelectedRow()[idColumnIndex]
				model := NewSecretModel(m.address, m.jwtCookie)
				cmd := getSecret(id, m.address, m.jwtCookie)
				return model, cmd
			}
		default:
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
	case listSecretsFailedMsg:
		{
			m.failtureStatusCode = msg.statusCode
			m.table.Blur()
		}
	case errMsg:
		{
			m.err = msg.err
			m.table.Blur()
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m secretsModel) View() string {
	s := strings.Builder{}

	if m.err != nil {
		s.WriteString(fmt.Sprintf(errTemplate, m.err.Error()))
	}
	if m.failtureStatusCode != 0 {
		s.WriteString(fmt.Sprintf(codeTemplate, m.failtureStatusCode))
	}

	s.WriteString(baseStyle.Render(m.table.View()) + "\n")

	return s.String()
}

func getSecret(id string, address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := newGetSecretCommand(id, address, jwtCookie)
	return cmd.execute
}
