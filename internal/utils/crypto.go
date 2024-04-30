package utils

import (
	"crypto/aes"
	"crypto/rand"

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
