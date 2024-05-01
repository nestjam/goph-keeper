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

	cipher, err := newBlockCipher(k.key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	nonce, _ := utils.GenerateRandom(cipher.NonceSize())

	sealed := dataKey.Copy()
	sealed.Key = cipher.Seal(nil, nonce, dataKey.Key, nil)
	sealed.IV = nonce

	return sealed, nil
}

func (k *MasterKey) Unseal(dataKey *DataKey) (*DataKey, error) {
	const op = "unseal"

	cipher, err := newBlockCipher(k.key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed := dataKey.Copy()
	unsealed.IV = nil
	unsealed.Key, err = open(dataKey.Key, dataKey.IV, cipher)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}
