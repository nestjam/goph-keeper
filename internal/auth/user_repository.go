package auth

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/auth/model"
)

var (
	ErrUserWithEmailIsRegistered  = errors.New("user with email has already been registered")
	ErrUserIsNotRegisteredAtEmail = errors.New("user is not registered at email")
)

type UserRepository interface {
	Register(user *model.User) (*model.User, error)
	GetByID(id uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
}
