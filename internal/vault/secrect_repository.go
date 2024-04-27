package vault

import (
	"context"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type SecretRepository interface {
	ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error)
	AddSecret(ctx context.Context, secret *model.Secret, userID uuid.UUID) (*model.Secret, error)
}
