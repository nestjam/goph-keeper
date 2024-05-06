package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type vaultServiceMock struct {
	ListSecretsFunc  func(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error)
	AddSecretFunc    func(ctx context.Context, secret *model.Secret, userID uuid.UUID) (uuid.UUID, error)
	UpdateSecretFunc func(ctx context.Context, secret *model.Secret, userID uuid.UUID) error
	GetSecretFunc    func(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error)
	DeleteSecretFunc func(ctx context.Context, secretID, userID uuid.UUID) error
}

func (m *vaultServiceMock) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	return m.ListSecretsFunc(ctx, userID)
}

func (m *vaultServiceMock) AddSecret(ctx context.Context, s *model.Secret, userID uuid.UUID) (uuid.UUID, error) {
	return m.AddSecretFunc(ctx, s, userID)
}

func (m *vaultServiceMock) UpdateSecret(ctx context.Context, s *model.Secret, userID uuid.UUID) error {
	return m.UpdateSecretFunc(ctx, s, userID)
}

func (m *vaultServiceMock) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	return m.GetSecretFunc(ctx, secretID, userID)
}

func (m *vaultServiceMock) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	return m.DeleteSecretFunc(ctx, secretID, userID)
}
