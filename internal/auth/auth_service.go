package auth

import "github.com/nestjam/goph-keeper/internal/auth/model"

type AuthService interface {
	Register(user *model.User) (*model.User, error)
}
