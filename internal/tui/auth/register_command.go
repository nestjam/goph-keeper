package auth

import (
	"crypto/tls"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	"github.com/nestjam/goph-keeper/internal/utils"
)

type registerCommand struct {
	address  string
	email    string
	password string
}

func newRegisterCommand(address, email, password string) registerCommand {
	return registerCommand{
		address:  address,
		email:    email,
		password: password,
	}
}

//nolint:dupl // loginCommand is not duplicate
func (c registerCommand) execute() tea.Msg {
	client := resty.New()

	//nolint:gosec // using self-signed certificate
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	req := httpAuth.RegisterUserRequest{
		Email:    c.email,
		Password: c.password,
	}
	url, err := url.JoinPath(c.address, "register")
	if err != nil {
		return errMsg{err}
	}

	resp, err := client.R().SetBody(req).Post(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		jwtCookie := findCookie(resp.Cookies(), utils.JWTCookieName)
		if jwtCookie == nil {
			return errMsg{ErrAuthTokenNotFound}
		}

		return registerCompletedMsg{
			jwtCookie: jwtCookie,
		}
	}

	return registerFailedMsg{resp.StatusCode()}
}
