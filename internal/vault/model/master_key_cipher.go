package model

import (
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/utils"
)

type MasterKeyCipher struct {
	masterKey *MasterKey
}

func NewMasterKeyCipher(masterKey *MasterKey) *MasterKeyCipher {
	return &MasterKeyCipher{masterKey}
}

func (c *MasterKeyCipher) Seal(dataKey *DataKey) (*DataKey, error) {
	const op = "seal"

	cipher := utils.NewBlockCipher(c.masterKey.key)
	ciphertext, err := cipher.Seal(dataKey.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	sealed := dataKey.Copy()
	sealed.Key = ciphertext

	return sealed, nil
}

func (c *MasterKeyCipher) Unseal(dataKey *DataKey) (*DataKey, error) {
	const op = "unseal"

	cipher := utils.NewBlockCipher(c.masterKey.key)
	plaintext, err := cipher.Unseal(dataKey.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed := dataKey.Copy()
	unsealed.Key = plaintext

	return unsealed, nil
}
