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
	}
}

func (k *DataKey) Seal(unsealed *Secret) (*Secret, error) {
	const op = "seal"

	cipher := utils.NewBlockCipher(k.Key)
	ciphertext, err := cipher.Seal(unsealed.Data)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	sealed := unsealed.Copy()
	sealed.Data = ciphertext
	sealed.KeyID = k.ID

	return sealed, nil
}

func (k *DataKey) Unseal(sealed *Secret) (unsealed *Secret, err error) {
	const op = "unseal"

	cipher := utils.NewBlockCipher(k.Key)
	plaintext, err := cipher.Unseal(sealed.Data)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed = sealed.Copy()
	unsealed.KeyID = uuid.Nil
	unsealed.Data = plaintext

	return unsealed, nil
}
