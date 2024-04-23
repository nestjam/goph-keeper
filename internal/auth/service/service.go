package service

import (
	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/pkg/errors"
)

type authService struct {
	repo auth.UserRepository
}

func NewAuthService(repo auth.UserRepository) auth.AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(user *model.User) (*model.User, error) {
	const op = "register user"

	user = user.Copy()
	err := user.HashPassword()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return createdUser, nil
}
