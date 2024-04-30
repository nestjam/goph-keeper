package inmemory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type dataKeyRepository struct {
	keys map[uuid.UUID]*model.DataKey
	key  *model.DataKey
	mu   sync.Mutex
}

func NewDataKeyRepository() vault.DataKeyRepository {
	return &dataKeyRepository{
		keys: make(map[uuid.UUID]*model.DataKey),
	}
}

func (r *dataKeyRepository) RotateKey(ctx context.Context, key *model.DataKey) (*model.DataKey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	newKey := key.Copy()
	id := uuid.New()
	newKey.ID = id

	r.keys[id] = newKey
	r.key = newKey

	return newKey, nil
}

func (r *dataKeyRepository) GetKey(ctx context.Context) (*model.DataKey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.key == nil {
		return nil, vault.ErrKeyNotFound
	}

	return r.key, nil
}

func (r *dataKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.DataKey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if key, ok := r.keys[id]; ok {
		return key, nil
	}

	return nil, vault.ErrKeyNotFound
}

func (r *dataKeyRepository) UpdateStats(ctx context.Context, key *model.DataKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, ok := r.keys[key.ID]
	if !ok {
		return vault.ErrKeyNotFound
	}

	stored.EncryptedSize = key.EncryptedSize

	return nil
}
