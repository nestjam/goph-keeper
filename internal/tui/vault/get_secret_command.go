package vault

import (
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type getSecretCommand struct {
	client    *resty.Client
	jwtCookie *http.Cookie
	address   string
	secretID  string
}

func newGetSecretCommand(secretID, addr string, jwt *http.Cookie, client *resty.Client) getSecretCommand {
	return getSecretCommand{
		secretID:  secretID,
		address:   addr,
		jwtCookie: jwt,
		client:    client,
	}
}

func (c getSecretCommand) execute() tea.Msg {
	url, err := url.JoinPath(c.address, baseURL, c.secretID)
	if err != nil {
		return getSecretFailedMsg{
			err:      err,
			secretID: c.secretID,
		}
	}
	var res httpVault.GetSecretResponse
	resp, err := c.client.R().SetResult(&res).SetCookie(c.jwtCookie).Get(url)
	if err != nil {
		return getSecretFailedMsg{
			err:      err,
			secretID: c.secretID,
		}
	}

	if resp.IsSuccess() {
		return getSecretCompletedMsg{res.Secret}
	}

	return getSecretFailedMsg{
		statusCode: resp.StatusCode(),
		secretID:   c.secretID,
	}
}
