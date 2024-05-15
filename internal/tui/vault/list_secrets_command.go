package vault

import (
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type listSecretsCommand struct {
	client    *resty.Client
	jwtCookie *http.Cookie
	address   string
}

func NewListSecretsCommand(addr string, jwt *http.Cookie, client *resty.Client) listSecretsCommand {
	return listSecretsCommand{
		address:   addr,
		jwtCookie: jwt,
		client:    client,
	}
}

func (c listSecretsCommand) Execute() tea.Msg {
	url, err := url.JoinPath(c.address, baseURL)
	if err != nil {
		return listSecretsFailedMsg{err: err}
	}

	var res struct {
		List []*httpVault.Secret `json:"list,omitempty"`
	}
	resp, err := c.client.R().SetResult(&res).SetCookie(c.jwtCookie).Get(url)
	if err != nil {
		return listSecretsFailedMsg{err: err}
	}

	if resp.IsSuccess() {
		return listSecretsCompletedMsg{res.List}
	}

	return listSecretsFailedMsg{statusCode: resp.StatusCode()}
}
