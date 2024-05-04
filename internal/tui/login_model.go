package tui

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginModel struct {
	serverAddress string
	email         string
	password      string
	err           error
	textinput     textinput.Model
}

func NewLoginModel() loginModel {
	ti := textinput.New()
	ti.Placeholder = "Enter server address"
	ti.Focus()

	return loginModel{
		textinput: ti,
	}
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := m.textinput.Value()
			acceptInput(&m, input)

			if canLogin(m.serverAddress, m.email, m.password) {
				return m, login(m.serverAddress, m.email, m.password)
			}
			return m, nil
		default:
		}
	case loginCompletedMsg:
		{
			return secretsModel{}, listSecrets(m.serverAddress, msg.jwtCookie)
		}
	case loginFailedMsg:
		{
			m.password = ""
			m.email = ""
			acceptServerAddress(&m, m.serverAddress)
		}
	case errMsg:
		{
			m.err = msg.err
			return m, nil
		}
	default:
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m loginModel) View() string {
	s := strings.Builder{}

	if m.serverAddress != "" {
		s.WriteString("server: ")
		s.WriteString(m.serverAddress + "\n\n")
	}
	if m.email != "" {
		s.WriteString("email: ")
		s.WriteString(m.email + "\n\n")
	}
	if m.password != "" {
		s.WriteString("password: ")
		s.WriteString(m.password + "\n\n")
	}
	if m.err != nil {
		s.WriteString("error: ")
		s.WriteString(m.err.Error() + "\n\n")
	}

	s.WriteString(m.textinput.View())

	return s.String()
}

func canLogin(address, email, password string) bool {
	return address != "" && email != "" && password != ""
}

func login(address, email, password string) tea.Cmd {
	cmd := loginCommand{
		address:  address,
		email:    email,
		password: password,
	}
	return cmd.execute
}

func listSecrets(address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := listSecretsCommand{
		address:   address,
		jwtCookie: jwtCookie,
	}
	return cmd.execute
}

func acceptInput(m *loginModel, input string) {
	if input == "" {
		return
	}

	switch {
	case m.serverAddress == "":
		acceptServerAddress(m, input)
	case m.email == "":
		acceptEmail(m, input)
	case m.password == "":
		acceptPassword(m, input)
	}
}

func acceptPassword(m *loginModel, input string) {
	m.password = input
	m.textinput.SetValue("")
	m.textinput.Placeholder = ""
}

func acceptEmail(m *loginModel, input string) {
	m.email = input
	m.textinput.SetValue("")
	m.textinput.Placeholder = "Enter password"
}

func acceptServerAddress(m *loginModel, input string) {
	m.serverAddress = input
	m.textinput.SetValue("")
	m.textinput.Placeholder = "Enter email"
}
