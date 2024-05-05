package vault

import (
	"crypto/tls"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type listSecretsCommand struct {
	jwtCookie *http.Cookie
	address   string
}

func NewListSecretsCommand(address string, jwtCookie *http.Cookie) listSecretsCommand {
	return listSecretsCommand{
		address:   address,
		jwtCookie: jwtCookie,
	}
}

func (c listSecretsCommand) Execute() tea.Msg {
	client := resty.New()

	//nolint:gosec // using self-signed certificate
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	url, err := url.JoinPath(c.address, baseURL)
	if err != nil {
		return errMsg{err}
	}
	var res httpVault.ListSecretsResponse
	resp, err := client.R().SetResult(&res).SetCookie(c.jwtCookie).Get(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		return listSecretsCompletedMsg{res.List}
	}

	return listSecretsFailedMsg{resp.StatusCode()}
}
