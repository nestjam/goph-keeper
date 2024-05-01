package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

const (
	encryptedDataSizeThreshold = 1 * 1024 * 1024 // 1GB
	encryptionsCountThreshold  = 1000
)

type KeyRotationConfig struct {
	EncryptedDataSizeThreshold int64
	EncryptionsCountThreshold  int
}

func NewKeyRotationConfig() KeyRotationConfig {
	return KeyRotationConfig{
		EncryptedDataSizeThreshold: encryptedDataSizeThreshold,
		EncryptionsCountThreshold:  encryptionsCountThreshold,
	}
}

type keyService struct {
	keyRepo vault.DataKeyRepository
	cipher  *model.Cipher
	config  KeyRotationConfig
}

func NewKeyService(keyRepo vault.DataKeyRepository, config KeyRotationConfig) *keyService {
	return &keyService{
		keyRepo: keyRepo,
		cipher:  model.NewAESGCMCipher(),
		config:  config,
	}
}

func (k *keyService) Seal(ctx context.Context, secret *model.Secret) (*model.Secret, error) {
	const op = "seal"

	key, err := k.keyRepo.GetKey(ctx)

	if errors.Is(err, vault.ErrKeyNotFound) ||
		key.EncryptedDataSize >= k.config.EncryptedDataSizeThreshold ||
		key.EncryptionsCount >= k.config.EncryptionsCountThreshold {
		newKey, err := model.NewDataKey()
		if err != nil {
			return nil, errors.Wrap(err, op)
		}

		key, err = k.keyRepo.RotateKey(ctx, newKey)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}
	}

	sealed, err := k.cipher.Seal(secret, key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	dataSize := int64(len(secret.Data))
	err = k.keyRepo.UpdateStats(ctx, key.ID, dataSize)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return sealed, nil
}

func (k *keyService) Unseal(ctx context.Context, secret *model.Secret) (*model.Secret, error) {
	const op = "unseal"

	key, err := k.keyRepo.GetByID(ctx, secret.KeyID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed, err := k.cipher.Unseal(secret, key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}
