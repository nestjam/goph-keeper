package vault

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type secretModel struct {
	textarea           textarea.Model
	err                error
	jwtCookie          *http.Cookie
	address            string
	secret             httpVault.Secret
	failtureStatusCode int
	isNew              bool
}

func NewSecretModel(address string, jwtCookie *http.Cookie) secretModel {
	ti := textarea.New()
	ti.Focus()

	return secretModel{
		textarea:  ti,
		address:   address,
		jwtCookie: jwtCookie,
	}
}

func (m secretModel) Init() tea.Cmd {
	return nil
}

func (m secretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case getSecretCompletedMsg:
		{
			m.secret = msg.secret
			m.textarea.SetValue(m.secret.Data)
		}
	case getSecretFailedMsg:
		{
			m.failtureStatusCode = msg.statusCode
			m.textarea.Blur()
		}
	case createSecretRequestedMsg:
		{
			m.secret = httpVault.Secret{}
			m.isNew = true
		}
	case saveSecretCompletedMsg:
		{
			m.secret = msg.secret
			m.textarea.SetValue(msg.secret.Data)
			m.isNew = false
		}
	case errMsg:
		{
			m.err = msg.err
			m.textarea.Blur()
		}
	default:
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m secretModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		model := NewSecretsModel(m.address, m.jwtCookie)
		cmd := listSecrets(m.address, m.jwtCookie)
		return model, cmd
	case tea.KeyCtrlS:
		secret := m.secret
		secret.Data = m.textarea.Value()
		cmd := saveSecret(secret, m.address, m.jwtCookie)
		return m, cmd
	default:
		{
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
	}
}

func (m secretModel) View() string {
	s := strings.Builder{}

	if m.err != nil {
		s.WriteString(fmt.Sprintf(errTemplate, m.err.Error()))
	}
	if m.failtureStatusCode != 0 {
		s.WriteString(fmt.Sprintf(codeTemplate, m.failtureStatusCode))
	}

	s.WriteString(fmt.Sprintf("id: %s\n\n", m.secret.ID))
	s.WriteString(m.textarea.View())

	return s.String()
}

func listSecrets(address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := NewListSecretsCommand(address, jwtCookie)
	return cmd.Execute
}

func saveSecret(secret httpVault.Secret, address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := newSaveSecretCommand(secret, address, jwtCookie)
	return cmd.execute
}
