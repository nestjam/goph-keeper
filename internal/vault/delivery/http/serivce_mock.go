package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type vaultServiceMock struct {
	ListSecretsFunc func(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error)
	AddSecretFunc   func(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error)
	GetSecretFunc   func(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error)
}

func (m *vaultServiceMock) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	return m.ListSecretsFunc(ctx, userID)
}

func (m *vaultServiceMock) AddSecret(ctx context.Context, s *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	return m.AddSecretFunc(ctx, s, userID)
}

func (m *vaultServiceMock) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	return m.GetSecretFunc(ctx, secretID, userID)
}
