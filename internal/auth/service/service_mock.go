package service

import (
	"context"

	"github.com/nestjam/goph-keeper/internal/auth/model"
)

type AuthServiceMock struct {
	RegisterFunc func(ctx context.Context, user *model.User) (*model.User, error)
}

func (s *AuthServiceMock) Register(ctx context.Context, user *model.User) (*model.User, error) {
	return s.RegisterFunc(ctx, user)
}
