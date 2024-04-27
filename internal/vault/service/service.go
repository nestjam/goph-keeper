package service

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type vaultService struct {
	repo vault.SecretRepository
}

func NewVaultService(repo vault.SecretRepository) vault.VaultService {
	return &vaultService{
		repo: repo,
	}
}

func (s *vaultService) ListSecrets(userID uuid.UUID) ([]*model.Secret, error) {
	const op = "list secrets"

	secrets, err := s.repo.ListSecrets(userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return secrets, nil
}

func (s *vaultService) AddSecret(secret *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	panic("unimplemented")
}
