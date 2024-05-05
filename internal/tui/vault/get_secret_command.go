package vault

import (
	"crypto/tls"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type getSecretCommand struct {
	jwtCookie *http.Cookie
	address   string
	secretID  string
}

func NewGetSecretCommand(secretID, address string, jwtCookie *http.Cookie) getSecretCommand {
	return getSecretCommand{
		secretID:  secretID,
		address:   address,
		jwtCookie: jwtCookie,
	}
}

func (c getSecretCommand) execute() tea.Msg {
	client := resty.New()

	//nolint:gosec // using self-signed certificate
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	url, err := url.JoinPath(c.address, baseURL, c.secretID)
	if err != nil {
		return errMsg{err}
	}
	var res httpVault.GetSecretResponse
	resp, err := client.R().SetResult(&res).SetCookie(c.jwtCookie).Get(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		return getSecretCompletedMsg{res.Secret}
	}

	return getSecretFailedMsg{resp.StatusCode()}
}

type getSecretCompletedMsg struct {
	secret httpVault.Secret
}

type getSecretFailedMsg struct {
	statusCode int
}
