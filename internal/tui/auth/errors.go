package auth

import "errors"

var (
	ErrAuthTokenNotFound = errors.New("auth cookie not found")
)
