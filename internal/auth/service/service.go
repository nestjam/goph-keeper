package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
)

type authService struct {
	repo auth.UserRepository
}

func NewAuthService(repo auth.UserRepository) auth.AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, user *model.User) (*model.User, error) {
	const op = "register user"

	user = user.Copy()
	err := user.HashPassword()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	newUser, err := s.repo.Register(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, user *model.User) (*model.User, error) {
	const op = "login"

	foundUser, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	if !foundUser.ComparePassword(user.Password) {
		return nil, auth.ErrInvalidPassword
	}

	return foundUser, nil
}
