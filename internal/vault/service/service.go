package service

import (
	"context"

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

func (s *vaultService) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	const op = "list secrets"

	secrets, err := s.repo.ListSecrets(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return secrets, nil
}

func (s *vaultService) AddSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	const op = "add secret"

	addedSecret, err := s.repo.AddSecret(ctx, secret, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return addedSecret, nil
}
