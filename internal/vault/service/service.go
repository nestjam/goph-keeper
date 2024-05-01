package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type vaultService struct {
	secretRepo vault.SecretRepository
	keyring    *keyService
}

func NewVaultService(secretRepo vault.SecretRepository,
	keyRepo vault.DataKeyRepository,
	rootKey *model.MasterKey) vault.VaultService {
	return &vaultService{
		secretRepo: secretRepo,
		keyring:    NewKeyService(keyRepo, NewKeyRotationConfig(), rootKey),
	}
}

func (s *vaultService) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	const op = "list secrets"

	secrets, err := s.secretRepo.ListSecrets(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return secrets, nil
}

func (s *vaultService) AddSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	const op = "add secret"

	sealed, err := s.keyring.Seal(ctx, secret)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	added, err := s.secretRepo.AddSecret(ctx, sealed, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return added, nil
}

func (s *vaultService) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	const op = "get secret"

	secret, err := s.secretRepo.GetSecret(ctx, secretID, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed, err := s.keyring.Unseal(ctx, secret)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}

func (s *vaultService) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	const op = "delete secret"

	err := s.secretRepo.DeleteSecret(ctx, secretID, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
