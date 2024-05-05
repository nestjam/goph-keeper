package auth

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nestjam/goph-keeper/internal/tui/vault"
)

var choices = []string{"login", "regiser"}

type loginModel struct {
	address   string
	email     string
	password  string
	err       error
	textinput textinput.Model
	cursor    int
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
		return handleKeyMsg(msg, m)
	case loginCompletedMsg:
		return vault.NewSecretsModel(m.address, msg.jwtCookie), listSecrets(m.address, msg.jwtCookie)
	case registerCompletedMsg:
		return vault.NewSecretsModel(m.address, msg.jwtCookie), listSecrets(m.address, msg.jwtCookie)
	case loginFailedMsg, registerFailedMsg:
		{
			m.password = ""
			m.email = ""
			acceptServerAddress(&m, m.address)
		}
	case errMsg:
		{
			m.err = msg.err
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func handleKeyMsg(msg tea.KeyMsg, m loginModel) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyEnter:
		input := m.textinput.Value()
		acceptInput(&m, input)

		if isValid(m.address, m.email, m.password) {
			if m.cursor == 0 {
				return m, login(m.address, m.email, m.password)
			}
			if m.cursor == 1 {
				return m, register(m.address, m.email, m.password)
			}
		}
		return m, nil
	case tea.KeyDown:
		m.cursor++
		if m.cursor >= len(choices) {
			m.cursor = 0
		}
	case tea.KeyUp:
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(choices) - 1
		}
	default:
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m loginModel) View() string {
	s := strings.Builder{}

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n")

	if m.address != "" {
		s.WriteString("server: ")
		s.WriteString(m.address)
		s.WriteString("\n")
	}
	if m.email != "" {
		s.WriteString("email: ")
		s.WriteString(m.email)
		s.WriteString("\n")
	}
	if m.password != "" {
		s.WriteString("password: ")
		s.WriteString(m.password)
		s.WriteString("\n")
	}
	if m.err != nil {
		s.WriteString("error: ")
		s.WriteString(m.err.Error())
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(m.textinput.View())

	return s.String()
}

func isValid(address, email, password string) bool {
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

func register(address, email, password string) tea.Cmd {
	cmd := newRegisterCommand(address, email, password)
	return cmd.execute
}

func listSecrets(address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := vault.NewListSecretsCommand(address, jwtCookie)
	return cmd.Execute
}

func acceptInput(m *loginModel, input string) {
	if input == "" {
		return
	}

	switch {
	case m.address == "":
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
	m.address = input
	m.textinput.SetValue("")
	m.textinput.Placeholder = "Enter email"
}
