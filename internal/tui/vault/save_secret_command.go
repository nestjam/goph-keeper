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
	jwtCookie *http.Cookie
	secret    httpVault.Secret
	address   string
}

func newSaveSecretCommand(secret httpVault.Secret, address string, jwtCookie *http.Cookie) saveSecretCommand {
	return saveSecretCommand{
		jwtCookie: jwtCookie,
		secret:    secret,
		address:   address,
	}
}

func (c saveSecretCommand) execute() tea.Msg {
	client := resty.New()

	//nolint:gosec // using self-signed certificate
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	url, err := url.JoinPath(c.address, baseURL)
	if err != nil {
		return errMsg{err}
	}

	req := httpVault.AddSecretRequest{
		Secret: c.secret,
	}
	var res httpVault.AddSecretResponse
	resp, err := client.R().SetBody(req).SetCookie(c.jwtCookie).SetResult(&res).Post(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		return saveSecretCompletedMsg{res.Secret}
	}

	return saveSecretFailedMsg{resp.StatusCode()}
}
