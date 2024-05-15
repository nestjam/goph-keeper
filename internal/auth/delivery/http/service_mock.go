package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/auth/model"
)

type authServiceMock struct {
	RegisterFunc func(ctx context.Context, user *model.User) (uuid.UUID, error)
	LoginFunc    func(ctx context.Context, user *model.User) (uuid.UUID, error)
}

func (s *authServiceMock) Register(ctx context.Context, user *model.User) (uuid.UUID, error) {
	return s.RegisterFunc(ctx, user)
}

func (s *authServiceMock) Login(ctx context.Context, user *model.User) (uuid.UUID, error) {
	return s.LoginFunc(ctx, user)
}
