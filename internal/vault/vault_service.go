package vault

import (
	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type VaultService interface {
	ListSecrets(userID uuid.UUID) ([]*model.Secret, error)
}
