package vault

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

var (
	ErrSecretDoesNotExist = errors.New("secret does not exist")
	ErrUserDoesNotExist   = errors.New("user does not exist")
)

type SecretRepository interface {
	ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error)
	AddSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error)
	GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error)
}
