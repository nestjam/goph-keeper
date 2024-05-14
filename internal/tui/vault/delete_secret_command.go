package vault

import (
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type deleteSecretCommand struct {
	client    *resty.Client
	jwtCookie *http.Cookie
	address   string
	secretID  string
}

func newDeleteSecretCommand(secretID, addr string, jwt *http.Cookie, client *resty.Client) deleteSecretCommand {
	return deleteSecretCommand{
		jwtCookie: jwt,
		address:   addr,
		secretID:  secretID,
		client:    client,
	}
}

func (c deleteSecretCommand) execute() tea.Msg {
	url, err := url.JoinPath(c.address, baseURL, c.secretID)
	if err != nil {
		return errMsg{err}
	}
	resp, err := c.client.R().SetCookie(c.jwtCookie).Delete(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		return deleteSecretCompletedMsg{c.secretID}
	}

	return deleteSecretFailedMsg{resp.StatusCode()}
}
