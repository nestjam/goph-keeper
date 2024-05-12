package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nestjam/goph-keeper/internal/tui/vault"
	"github.com/nestjam/goph-keeper/internal/tui/vault/cache"
)

var choices = []string{"login", "regiser"}

type loginKeyMap struct {
	Quit     key.Binding
	Continue key.Binding
	Up       key.Binding
	Down     key.Binding
}

func (k loginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Continue, k.Quit}
}

func (k loginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

type loginModel struct {
	err          error
	help         help.Model
	address      string
	email        string
	password     string
	BuildVersion string
	BuildDate    string
	keys         loginKeyMap
	textinput    textinput.Model
	cursor       int
}

func NewLoginModel() loginModel {
	ti := textinput.New()
	ti.Placeholder = "Enter server address"
	ti.Focus()

	keys := loginKeyMap{
		Quit: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String(), tea.KeyCtrlC.String()),
			key.WithHelp("ctr+c", "quit"),
		),
		Continue: key.NewBinding(
			key.WithKeys(tea.KeyEnter.String()),
			key.WithHelp("enter", "continue"),
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

	return loginModel{
		keys:      keys,
		help:      help.New(),
		textinput: ti,
	}
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		return handleKeyMsg(msg, m)
	case loginCompletedMsg:
		cache := cache.New()
		return vault.NewSecretsModel(m.address, msg.jwtCookie, cache), listSecrets(m.address, msg.jwtCookie)
	case registerCompletedMsg:
		cache := cache.New()
		return vault.NewSecretsModel(m.address, msg.jwtCookie, cache), listSecrets(m.address, msg.jwtCookie)
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
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Continue):
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
	case key.Matches(msg, m.keys.Down):
		m.cursor++
		if m.cursor >= len(choices) {
			m.cursor = 0
		}
	case key.Matches(msg, m.keys.Up):
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

	s.WriteString(fmt.Sprintf("ver: %s, date: %s\n", m.BuildVersion, m.BuildDate))
	s.WriteString("\n")

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

	s.WriteString(m.textinput.View())

	s.WriteString("\n\n")
	s.WriteString(m.help.View(m.keys))

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
