package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type vaultService struct {
	repo   vault.SecretRepository
	cipher *model.Cipher
	key    []byte
}

func NewVaultService(repo vault.SecretRepository,
	rootKey []byte) vault.VaultService {
	return &vaultService{
		repo:   repo,
		cipher: model.NewAESGCMCipher(),
		key:    rootKey,
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

	sealed, err := s.cipher.Seal(secret, s.key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	added, err := s.repo.AddSecret(ctx, sealed, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return added, nil
}

func (s *vaultService) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	const op = "get secret"

	secret, err := s.repo.GetSecret(ctx, secretID, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed, err := s.cipher.Unseal(secret, s.key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}

func (s *vaultService) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	const op = "delete secret"

	err := s.repo.DeleteSecret(ctx, secretID, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
