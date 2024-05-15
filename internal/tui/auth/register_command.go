package auth

import (
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	"github.com/nestjam/goph-keeper/internal/utils"
)

type registerCommand struct {
	client   *resty.Client
	address  string
	email    string
	password string
}

func newRegisterCommand(address, email, password string, client *resty.Client) registerCommand {
	return registerCommand{
		address:  address,
		email:    email,
		password: password,
		client:   client,
	}
}

//nolint:dupl // loginCommand is not duplicate
func (c registerCommand) execute() tea.Msg {
	req := httpAuth.RegisterUserRequest{
		Email:    c.email,
		Password: c.password,
	}
	url, err := url.JoinPath(c.address, "register")
	if err != nil {
		return errMsg{err}
	}

	resp, err := c.client.R().SetBody(req).Post(url)
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
