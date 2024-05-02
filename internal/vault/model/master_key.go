package model

import (
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/utils"
)

type MasterKey struct {
	key []byte
}

func NewMasterKey(key []byte) *MasterKey {
	return &MasterKey{
		key: key,
	}
}

func (k *MasterKey) Seal(dataKey *DataKey) (*DataKey, error) {
	const op = "seal"

	cipher := utils.NewBlockCipher(k.key)
	ciphertext, err := cipher.Seal(dataKey.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	sealed := dataKey.Copy()
	sealed.Key = ciphertext

	return sealed, nil
}

func (k *MasterKey) Unseal(dataKey *DataKey) (*DataKey, error) {
	const op = "unseal"

	cipher := utils.NewBlockCipher(k.key)
	plaintext, err := cipher.Unseal(dataKey.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed := dataKey.Copy()
	unsealed.Key = plaintext

	return unsealed, nil
}
