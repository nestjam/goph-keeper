package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/utils"
)

type DataKeyCipher struct {
	dataKey *DataKey
}

func NewDataKeyCipher(dataKey *DataKey) *DataKeyCipher {
	return &DataKeyCipher{dataKey}
}

func (c *DataKeyCipher) Seal(unsealed *Secret) (*Secret, error) {
	const op = "seal"

	cipher := utils.NewBlockCipher(c.dataKey.Key)
	ciphertext, err := cipher.Seal(unsealed.Data)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	sealed := unsealed.Copy()
	sealed.Data = ciphertext
	sealed.KeyID = c.dataKey.ID

	return sealed, nil
}

func (c *DataKeyCipher) Unseal(sealed *Secret) (unsealed *Secret, err error) {
	const op = "unseal"

	cipher := utils.NewBlockCipher(c.dataKey.Key)
	plaintext, err := cipher.Unseal(sealed.Data)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed = sealed.Copy()
	unsealed.KeyID = uuid.Nil
	unsealed.Data = plaintext

	return unsealed, nil
}
