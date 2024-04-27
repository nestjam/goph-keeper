package http

import (
	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type vaultServiceMock struct {
	ListSecretsFunc func(userID uuid.UUID) ([]*model.Secret, error)
}

func (s *vaultServiceMock) ListSecrets(userID uuid.UUID) ([]*model.Secret, error) {
	return s.ListSecretsFunc(userID)
}
