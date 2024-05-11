package service

import (
	"context"

	"github.com/google/uuid"
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

func (s *authService) Register(ctx context.Context, user *model.User) (uuid.UUID, error) {
	const op = "register user"

	err := user.HashPassword()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	userID, err := s.repo.Register(ctx, user)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	return userID, nil
}

func (s *authService) Login(ctx context.Context, user *model.User) (uuid.UUID, error) {
	const op = "login"

	foundUser, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	if !foundUser.ComparePassword(user.Password) {
		return uuid.Nil, auth.ErrInvalidPassword
	}

	return foundUser.ID, nil
}
