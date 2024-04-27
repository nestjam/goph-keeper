package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type vaultServiceMock struct {
	ListSecretsFunc func(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error)
}

func (s *vaultServiceMock) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	return s.ListSecretsFunc(ctx, userID)
}
