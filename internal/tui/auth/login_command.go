package auth

import (
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	"github.com/nestjam/goph-keeper/internal/utils"
)

type loginCommand struct {
	client   *resty.Client
	address  string
	email    string
	password string
}

func newLoginCommand(address, email, password string, client *resty.Client) loginCommand {
	return loginCommand{
		address:  address,
		email:    email,
		password: password,
		client:   client,
	}
}

//nolint:dupl // registerCommand is not duplicate
func (c loginCommand) execute() tea.Msg {
	req := httpAuth.LoginUserRequest{
		Email:    c.email,
		Password: c.password,
	}
	url, err := url.JoinPath(c.address, "login")
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

		return loginCompletedMsg{
			jwtCookie: jwtCookie,
		}
	}

	return loginFailedMsg{resp.StatusCode()}
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == name {
			return cookies[i]
		}
	}

	return nil
}
