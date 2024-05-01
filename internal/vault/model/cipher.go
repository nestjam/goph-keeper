package model

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/utils"
)

type Cipher struct {
}

func NewAESGCMCipher() *Cipher {
	return &Cipher{}
}

func (c *Cipher) Seal(unsealed *Secret, dataKey *DataKey) (*Secret, error) {
	const op = "seal"

	cipher, err := newBlockCipher(dataKey.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	nonce, _ := utils.GenerateRandom(cipher.NonceSize())

	sealed := unsealed.Copy()
	sealed.Data = cipher.Seal(nil, nonce, unsealed.Data, nil)
	sealed.IV = nonce
	sealed.KeyID = dataKey.ID

	return sealed, nil
}

func (c *Cipher) Unseal(sealed *Secret, dataKey *DataKey) (unsealed *Secret, err error) {
	const op = "unseal"

	cipher, err := newBlockCipher(dataKey.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed = sealed.Copy()
	unsealed.IV = nil
	unsealed.KeyID = uuid.Nil
	unsealed.Data, err = open(sealed.Data, sealed.IV, cipher)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}

func open(sealed, iv []byte, cipher cipher.AEAD) (unsealed []byte, err error) {
	const op = "open"

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s: %v", op, r)
		}
	}()

	unsealed, err = cipher.Open(nil, iv, sealed, nil)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}

func newBlockCipher(key []byte) (cipher.AEAD, error) {
	const op = "new block cipher"

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return aead, nil
}
