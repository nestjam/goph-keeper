package auth

import (
	"context"
	"errors"

	"github.com/nestjam/goph-keeper/internal/auth/model"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type AuthService interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, user *model.User) (*model.User, error)
}
