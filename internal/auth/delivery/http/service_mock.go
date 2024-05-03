package http

import (
	"context"

	"github.com/nestjam/goph-keeper/internal/auth/model"
)

type authServiceMock struct {
	RegisterFunc func(ctx context.Context, user model.User) (model.User, error)
	LoginFunc    func(ctx context.Context, user model.User) (model.User, error)
}

func (s *authServiceMock) Register(ctx context.Context, user model.User) (model.User, error) {
	return s.RegisterFunc(ctx, user)
}

func (s *authServiceMock) Login(ctx context.Context, user model.User) (model.User, error) {
	return s.LoginFunc(ctx, user)
}
