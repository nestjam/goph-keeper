package auth

import (
	"context"

	"github.com/nestjam/goph-keeper/internal/auth/model"
)

type AuthService interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
}
