package model

import (
	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/pkg/errors"
)

type DataKey struct {
	Key               []byte
	ID                uuid.UUID
	EncryptedDataSize int64
	EncryptionsCount  int
}

func NewDataKey() (*DataKey, error) {
	key, err := utils.GenerateRandomAES256Key()
	if err != nil {
		return nil, errors.Wrap(err, "new aes-256 data key")
	}

	return &DataKey{Key: key}, nil
}

func (k *DataKey) Copy() *DataKey {
	return &DataKey{
		ID:                k.ID,
		Key:               k.Key,
		EncryptedDataSize: k.EncryptedDataSize,
		EncryptionsCount:  k.EncryptionsCount,
	}
}
