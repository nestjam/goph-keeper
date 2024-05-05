package vault

import (
	"crypto/tls"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type deleteSecretCommand struct {
	jwtCookie *http.Cookie
	address   string
	secretID  string
}

func newDeleteSecretCommand(secretID, address string, jwtCookie *http.Cookie) deleteSecretCommand {
	return deleteSecretCommand{
		jwtCookie: jwtCookie,
		address:   address,
		secretID:  secretID,
	}
}

func (c deleteSecretCommand) execute() tea.Msg {
	client := resty.New()

	//nolint:gosec // using self-signed certificate
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	url, err := url.JoinPath(c.address, baseURL, c.secretID)
	if err != nil {
		return errMsg{err}
	}
	resp, err := client.R().SetCookie(c.jwtCookie).Delete(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		return deleteSecretCompletedMsg{c.secretID}
	}

	return deleteSecretFailedMsg{resp.StatusCode()}
}
