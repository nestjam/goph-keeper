package vault

import (
	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type SecretRepository interface {
	ListSecrets(userID uuid.UUID) ([]*model.Secret, error)
	AddSecret(secret *model.Secret, userID uuid.UUID) (*model.Secret, error)
}
