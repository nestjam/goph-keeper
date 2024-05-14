package vault

import (
	"crypto/tls"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type saveSecretCommand struct {
	client    *resty.Client
	jwtCookie *http.Cookie
	secret    httpVault.Secret
	address   string
}

func newSaveSecretCommand(secret httpVault.Secret, addr string, jwt *http.Cookie, c *resty.Client) saveSecretCommand {
	return saveSecretCommand{
		jwtCookie: jwt,
		secret:    secret,
		address:   addr,
		client:    c,
	}
}

func (c saveSecretCommand) execute() tea.Msg {
	if c.secret.ID == "" {
		return c.addSecret()
	}

	return c.updateSecret()
}

func (c saveSecretCommand) addSecret() tea.Msg {
	url, err := url.JoinPath(c.address, baseURL)
	if err != nil {
		return errMsg{err}
	}

	req := httpVault.AddSecretRequest{
		Secret: c.secret,
	}
	var res httpVault.AddSecretResponse
	resp, err := c.client.R().SetBody(req).SetCookie(c.jwtCookie).SetResult(&res).Post(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		secret := res.Secret
		secret.Data = c.secret.Data
		return saveSecretCompletedMsg{secret}
	}

	return saveSecretFailedMsg{resp.StatusCode()}
}

func (c saveSecretCommand) updateSecret() tea.Msg {
	client := resty.New()

	//nolint:gosec // using self-signed certificate
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	url, err := url.JoinPath(c.address, baseURL, c.secret.ID)
	if err != nil {
		return errMsg{err}
	}

	req := httpVault.UpdateSecretRequest{
		Secret: c.secret,
	}
	resp, err := client.R().SetBody(req).SetCookie(c.jwtCookie).Patch(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		return saveSecretCompletedMsg{c.secret}
	}

	return saveSecretFailedMsg{resp.StatusCode()}
}
