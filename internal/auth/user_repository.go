package auth

import (
	"context"
	"errors"

	"github.com/nestjam/goph-keeper/internal/auth/model"
)

var (
	ErrUserWithEmailIsRegistered = errors.New("user with email has already been registered")
	ErrUserPasswordIsEmpty       = errors.New("user password is empty")
	ErrUserIsNotRegistered       = errors.New("user is not registered")
)

type UserRepository interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
}
