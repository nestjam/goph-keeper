package vault

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nestjam/goph-keeper/internal/tui/vault/cache"
	vault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type secretKeyMap struct {
	Quit   key.Binding
	Save   key.Binding
	Return key.Binding
}

func (k secretKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Return, k.Quit}
}

func (k secretKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

type secretModel struct {
	textarea           textarea.Model
	err                error
	jwtCookie          *http.Cookie
	cache              *cache.SecretsCache
	help               help.Model
	secret             vault.Secret
	address            string
	keys               secretKeyMap
	failtureStatusCode int
	isNew              bool
	isOffline          bool
}

func NewSecretModel(address string, jwtCookie *http.Cookie, cache *cache.SecretsCache) secretModel {
	ti := textarea.New()
	ti.Focus()

	keys := secretKeyMap{
		Quit: key.NewBinding(
			key.WithKeys(tea.KeyCtrlC.String()),
			key.WithHelp("ctrl+c", quitApp),
		),
		Save: key.NewBinding(
			key.WithKeys(tea.KeyCtrlS.String()),
			key.WithHelp("ctrl+s", "save"),
		),
		Return: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String()),
			key.WithHelp("esc", "return"),
		),
	}

	return secretModel{
		keys:      keys,
		help:      help.New(),
		textarea:  ti,
		address:   address,
		jwtCookie: jwtCookie,
		cache:     cache,
	}
}

func (m secretModel) Init() tea.Cmd {
	return nil
}

func (m secretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case getSecretCompletedMsg:
		{
			m.secret = msg.secret
			m.cache.CacheSecret(&msg.secret)
			m.textarea.SetValue(m.secret.Data)
		}
	case getSecretFailedMsg:
		{
			m.failtureStatusCode = msg.statusCode
			m.setOfflineMode(true)
			m.textarea.Blur()

			if secret, ok := m.cache.GetSecret(msg.secretID); ok {
				m.textarea.SetValue(secret.Data)
			}
		}
	case createSecretRequestedMsg:
		{
			m.secret = vault.Secret{}
			m.isNew = true
		}
	case saveSecretCompletedMsg:
		{
			m.secret = msg.secret
			m.cache.CacheSecret(&msg.secret)
			m.textarea.SetValue(msg.secret.Data)
			m.isNew = false
		}
	case errMsg:
		{
			m.err = msg.err
			m.setOfflineMode(true)
			m.textarea.Blur()
		}
	default:
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m secretModel) View() string {
	s := strings.Builder{}

	if m.isOffline {
		s.WriteString("offline mode")
		s.WriteString("\n")
	}

	if m.err != nil {
		s.WriteString(fmt.Sprintf(errTemplate, m.err.Error()))
	}
	if m.failtureStatusCode != 0 {
		s.WriteString(fmt.Sprintf(codeTemplate, m.failtureStatusCode))
	}

	s.WriteString(fmt.Sprintf("id: %s\n\n", m.secret.ID))
	s.WriteString(m.textarea.View())

	s.WriteString("\n\n")
	s.WriteString(m.help.View(m.keys))

	return s.String()
}

func (m *secretModel) setOfflineMode(v bool) {
	m.isOffline = v
	m.keys.Save.SetEnabled(!v)
}

func (m secretModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Return):
		model := NewSecretsModel(m.address, m.jwtCookie, m.cache)
		cmd := listSecrets(m.address, m.jwtCookie)
		return model, cmd
	case key.Matches(msg, m.keys.Save):
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

func listSecrets(address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := NewListSecretsCommand(address, jwtCookie)
	return cmd.Execute
}

func saveSecret(secret vault.Secret, address string, jwtCookie *http.Cookie) tea.Cmd {
	cmd := newSaveSecretCommand(secret, address, jwtCookie)
	return cmd.execute
}
