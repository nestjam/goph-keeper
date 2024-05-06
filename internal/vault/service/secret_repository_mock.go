package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type secretRepositoryMock struct {
	ListSecretsFunc  func(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error)
	AddSecretFunc    func(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error)
	UpdateSecretFunc func(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error)
	GetSecretFunc    func(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error)
	DeleteSecretFunc func(ctx context.Context, secretID, userID uuid.UUID) error
}

func (m *secretRepositoryMock) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	return m.ListSecretsFunc(ctx, userID)
}

func (m *secretRepositoryMock) AddSecret(ctx context.Context, s *model.Secret, u uuid.UUID) (*model.Secret, error) {
	return m.AddSecretFunc(ctx, s, u)
}

func (m *secretRepositoryMock) UpdateSecret(ctx context.Context, s *model.Secret, u uuid.UUID) (*model.Secret, error) {
	return m.UpdateSecretFunc(ctx, s, u)
}

func (m *secretRepositoryMock) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	return m.GetSecretFunc(ctx, secretID, userID)
}

func (m *secretRepositoryMock) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	return m.DeleteSecretFunc(ctx, secretID, userID)
}
