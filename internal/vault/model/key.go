package model

import (
	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/pkg/errors"
)

type DataKey struct {
	Key           []byte
	ID            uuid.UUID
	EncryptedSize int64
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
		ID:            k.ID,
		Key:           k.Key,
		EncryptedSize: k.EncryptedSize,
	}
}
