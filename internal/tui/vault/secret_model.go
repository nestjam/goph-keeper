package vault

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

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
	client             *resty.Client
	jwtCookie          *http.Cookie
	cache              *cache.SecretsCache
	help               help.Model
	secret             vault.Secret
	address            string
	keys               secretKeyMap
	failtureStatusCode int
	isNew              bool
	isOffline          bool
	dataCached         bool
}

func NewSecretModel(addr string, jwt *http.Cookie, cache *cache.SecretsCache, client *resty.Client) secretModel {
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
		address:   addr,
		jwtCookie: jwt,
		cache:     cache,
		client:    client,
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
			m.err = msg.err
			m.failtureStatusCode = msg.statusCode
			m.setOfflineMode(true)

			if secret, dataCached, ok := m.cache.GetSecret(msg.secretID); ok {
				m.secret = *secret
				m.dataCached = dataCached
				if dataCached {
					m.textarea.SetValue(secret.Data)
				} else {
					m.textarea.Placeholder = noCachedData
				}
			}
			m.textarea.Blur()
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
		s.WriteString(offlineMode)
		s.WriteString("\n")
	}

	if m.err != nil {
		s.WriteString(fmt.Sprintf(errTemplate, m.err.Error()))
	}
	if m.failtureStatusCode != 0 {
		s.WriteString(fmt.Sprintf(codeTemplate, m.failtureStatusCode))
	}

	s.WriteString(fmt.Sprintf("id: %s", m.secret.ID))
	s.WriteString("\n\n")

	s.WriteString(m.textarea.View())

	// hot keys help
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
		model := NewSecretsModel(m.address, m.jwtCookie, m.cache, m.client)
		cmd := listSecrets(m.address, m.jwtCookie, m.client)
		return model, cmd
	case key.Matches(msg, m.keys.Save):
		secret := m.secret
		secret.Data = m.textarea.Value()
		cmd := saveSecret(secret, m.address, m.jwtCookie, m.client)
		return m, cmd
	default:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}
}

func listSecrets(addr string, jwt *http.Cookie, client *resty.Client) tea.Cmd {
	cmd := NewListSecretsCommand(addr, jwt, client)
	return cmd.Execute
}

func saveSecret(secret vault.Secret, addr string, jwt *http.Cookie, client *resty.Client) tea.Cmd {
	cmd := newSaveSecretCommand(secret, addr, jwt, client)
	return cmd.execute
}
