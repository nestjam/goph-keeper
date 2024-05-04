package tui

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"

	httpAuth "github.com/nestjam/goph-keeper/internal/auth/delivery/http"
	"github.com/nestjam/goph-keeper/internal/utils"
)

type loginCommand struct {
	address  string
	email    string
	password string
}

func (c loginCommand) execute() tea.Msg {
	client := resty.New()

	//nolint:gosec // using self-signed certificate
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	req := httpAuth.LoginUserRequest{
		Email:    c.email,
		Password: c.password,
	}
	url, err := url.JoinPath(c.address, "login")
	if err != nil {
		return errMsg{err}
	}
	resp, err := client.R().SetBody(req).Post(url)
	if err != nil {
		return errMsg{err}
	}

	if resp.IsSuccess() {
		jwtCookie := findCookie(resp, utils.JWTCookieName)
		if jwtCookie == nil {
			return errMsg{errors.New("auth cookie not found")}
		}

		return loginCompletedMsg{
			jwtCookie: jwtCookie,
		}
	}

	return loginFailedMsg{resp.StatusCode()}
}

func findCookie(r *resty.Response, name string) *http.Cookie {
	cookies := r.Cookies()

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == name {
			return cookies[i]
		}
	}

	return nil
}

type loginCompletedMsg struct {
	jwtCookie *http.Cookie
}

type errMsg struct {
	err error
}

type loginFailedMsg struct {
	statusCode int
}
