package vault

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

var (
	ErrSecretNotFound = errors.New("secret not found")
)

//nolint:dupl // SecretRepository is not duplicate of VaultService
type SecretRepository interface {
	ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error)
	AddSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error)
	UpdateSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) error
	GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error)
	DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error
}
