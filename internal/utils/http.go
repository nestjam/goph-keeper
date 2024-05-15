package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/config"
)

const (
	JWTCookieName = "jwt"
	UserIDClaim   = "user_id"
	JWTAlg        = "HS256"
)

var (
	ErrUserIDNotFound = errors.New("user id not found")
	ErrInvalidUserID  = errors.New("invalid user id")
)

type AuthCookieBaker struct {
	auth   *jwtauth.JWTAuth
	config config.JWTAuthConfig
}

func NewAuthCookieBaker(config config.JWTAuthConfig) *AuthCookieBaker {
	return &AuthCookieBaker{
		auth:   jwtauth.New(JWTAlg, []byte(config.SignKey), nil),
		config: config,
	}
}

func (h *AuthCookieBaker) JWTAuth() *jwtauth.JWTAuth {
	return h.auth
}

func (h *AuthCookieBaker) BakeCookie(userID uuid.UUID) (*http.Cookie, error) {
	const op = "bake cookie"

	claims := make(map[string]interface{})
	claims[UserIDClaim] = userID.String()
	jwtauth.SetExpiryIn(claims, h.config.TokenExpiryIn)

	_, token, err := h.auth.Encode(claims)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	cookie := &http.Cookie{
		Name:     JWTCookieName,
		Value:    token,
		MaxAge:   int(h.config.TokenExpiryIn / time.Second),
		HttpOnly: true,
	}
	return cookie, nil
}

func UserFromContext(ctx context.Context) (uuid.UUID, error) {
	const op = "user from context"
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	claim, ok := claims[UserIDClaim]
	if !ok {
		return uuid.Nil, errors.Wrap(ErrUserIDNotFound, op)
	}

	s, _ := claim.(string)
	userID, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	return userID, nil
}
