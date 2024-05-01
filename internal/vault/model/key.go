package model

import (
	"crypto/aes"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/utils"
)

const DataKeySize = 2 * aes.BlockSize

type DataKey struct {
	Key               []byte
	IV                []byte
	ID                uuid.UUID
	EncryptedDataSize int64
	EncryptionsCount  int
}

func NewDataKey() (*DataKey, error) {
	key, err := utils.GenerateRandom(DataKeySize)
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
		IV:                k.IV,
	}
}
