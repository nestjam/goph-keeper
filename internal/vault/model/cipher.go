package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/utils"
)

type SecretCipher struct {
}

func NewSecretCipher() *SecretCipher {
	return &SecretCipher{}
}

func (c *SecretCipher) Seal(unsealed *Secret, dataKey *DataKey) (*Secret, error) {
	const op = "seal"

	cipher := utils.NewBlockCipher(dataKey.Key)
	ciphertext, iv, err := cipher.Seal(unsealed.Data)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	sealed := unsealed.Copy()
	sealed.Data = ciphertext
	sealed.IV = iv
	sealed.KeyID = dataKey.ID

	return sealed, nil
}

func (c *SecretCipher) Unseal(sealed *Secret, dataKey *DataKey) (unsealed *Secret, err error) {
	const op = "unseal"

	cipher := utils.NewBlockCipher(dataKey.Key)
	plaintext, err := cipher.Unseal(sealed.Data, sealed.IV)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed = sealed.Copy()
	unsealed.IV = nil
	unsealed.KeyID = uuid.Nil
	unsealed.Data = plaintext

	return unsealed, nil
}
