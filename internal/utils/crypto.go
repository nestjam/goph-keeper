package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/pkg/errors"
)

func GenerateRandom(size int) ([]byte, error) {
	const op = "generate random"

	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return b, nil
}

func GenerateRandomAES256Key() ([]byte, error) {
	const (
		op      = "generate random aes-256 key"
		keySize = 2 * aes.BlockSize
	)

	b, err := GenerateRandom(keySize)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return b, nil
}

type blockCipher struct {
	key []byte
}

func NewBlockCipher(key []byte) *blockCipher {
	return &blockCipher{key: key}
}

func (c *blockCipher) Seal(plaintext []byte) (ciphertext []byte, err error) {
	const op = "seal"

	cipher, err := newBlockCipher(c.key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	nonce, _ := GenerateRandom(cipher.NonceSize())

	return cipher.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *blockCipher) Unseal(ciphertext []byte) (plaintext []byte, err error) {
	const op = "seal"

	cipher, err := newBlockCipher(c.key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	plaintext, err = open(ciphertext, cipher)
	return
}

func open(ciphertext []byte, cipher cipher.AEAD) (plaintext []byte, err error) {
	const op = "open"

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s: %v", op, r)
		}
	}()

	plaintext, err = cipher.Open(nil, ciphertext[:cipher.NonceSize()], ciphertext[cipher.NonceSize():], nil)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return plaintext, nil
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
