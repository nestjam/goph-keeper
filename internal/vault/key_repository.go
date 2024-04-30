package vault

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

var (
	ErrKeyNotFound = errors.New("key not fount")
)

type DataKeyRepository interface {
	RotateKey(ctx context.Context, key *model.DataKey) (*model.DataKey, error)
	GetKey(ctx context.Context) (*model.DataKey, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.DataKey, error)
	UpdateStats(ctx context.Context, key *model.DataKey) error
}
