package vault

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

const idColumnIndex = 1

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type secretsKeyMap struct {
	Quit   key.Binding
	Up     key.Binding
	Down   key.Binding
	Edit   key.Binding
	Delete key.Binding
	Add    key.Binding
}

func (k secretsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Edit, k.Delete, k.Quit}
}

func (k secretsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

type secretsModel struct {
	keys               secretsKeyMap
	help               help.Model
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
		tableHeight = 10
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

	keys := secretsKeyMap{
		Quit: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String(), tea.KeyCtrlC.String()),
			key.WithHelp("esc", "quit"),
		),
		Add: key.NewBinding(
			key.WithKeys(tea.KeyCtrlN.String()),
			key.WithHelp("ctrl+n", "create"),
		),
		Edit: key.NewBinding(
			key.WithKeys(tea.KeyEnter.String()),
			key.WithHelp("enter", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys(tea.KeyDelete.String()),
			key.WithHelp("del", "delete"),
		),
		Up: key.NewBinding(
			key.WithKeys(tea.KeyUp.String()),
			key.WithHelp("↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys(tea.KeyDown.String()),
			key.WithHelp("↓", "move down"),
		),
	}

	return secretsModel{
		keys:      keys,
		help:      help.New(),
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
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
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
	case deleteSecretCompletedMsg:
		{
			rows := m.table.Rows()
			rows = deleteRow(msg.secretID, rows)
			m.table.SetRows(rows)
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

func (m secretsModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Edit):
		{
			id := m.getSelectedSecretID()
			model := NewSecretModel(m.address, m.jwtCookie)
			cmd := getSecret(id, m.address, m.jwtCookie)
			return model, cmd
		}
	case key.Matches(msg, m.keys.Delete):
		{
			id := m.getSelectedSecretID()
			model := m
			cmd := deleteSecret(id, m.address, m.jwtCookie)
			return model, cmd
		}
	case key.Matches(msg, m.keys.Add):
		{
			model := NewSecretModel(m.address, m.jwtCookie)
			cmd := createSecret()
			return model, cmd
		}
	default:
		{
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}
	}
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

	s.WriteString("\n")
	s.WriteString(m.help.View(m.keys))

	return s.String()
}

func (m *secretsModel) getSelectedSecretID() string {
	return m.table.SelectedRow()[idColumnIndex]
}

func getSecret(id string, address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := newGetSecretCommand(id, address, jwtCookie)
	return cmd.execute
}

func deleteSecret(id string, address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := newDeleteSecretCommand(id, address, jwtCookie)
	return cmd.execute
}

func createSecret() tea.Cmd {
	cmd := newCreateSecretCommand()
	return cmd.execute
}

func deleteRow(id string, rows []table.Row) []table.Row {
	i := findIndex(id, rows)
	if i < 0 {
		return rows
	}
	return append(rows[:i], rows[i+1:]...)
}

func findIndex(id string, rows []table.Row) int {
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if row[idColumnIndex] == id {
			return i
		}
	}
	return -1
}
