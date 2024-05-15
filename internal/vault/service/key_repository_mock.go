package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type keyRepositoryMock struct {
	RotateKeyFunc   func(ctx context.Context, key *model.DataKey) (*model.DataKey, error)
	GetKeyFunc      func(ctx context.Context) (*model.DataKey, error)
	GetByIDFunc     func(ctx context.Context, id uuid.UUID) (*model.DataKey, error)
	UpdateStatsFunc func(ctx context.Context, id uuid.UUID, dataSize int64) error
}

func (m *keyRepositoryMock) RotateKey(ctx context.Context, key *model.DataKey) (*model.DataKey, error) {
	return m.RotateKeyFunc(ctx, key)
}

func (m *keyRepositoryMock) GetKey(ctx context.Context) (*model.DataKey, error) {
	return m.GetKeyFunc(ctx)
}

func (m *keyRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*model.DataKey, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *keyRepositoryMock) UpdateStats(ctx context.Context, id uuid.UUID, dataSize int64) error {
	return m.UpdateStatsFunc(ctx, id, dataSize)
}
