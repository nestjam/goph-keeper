package auth

import "net/http"

type registerCompletedMsg struct {
	jwtCookie *http.Cookie
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

type registerFailedMsg struct {
	statusCode int
}
