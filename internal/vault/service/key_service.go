package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

const (
	encryptedSizeThreshold = 1 * 1024 * 1024 // 1GB
)

type KeyRotationConfig struct {
	EncryptedSizeThreshold int64
}

func NewKeyRotationConfig() KeyRotationConfig {
	return KeyRotationConfig{
		EncryptedSizeThreshold: encryptedSizeThreshold,
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

	if errors.Is(err, vault.ErrKeyNotFound) || key.EncryptedSize >= k.config.EncryptedSizeThreshold {
		newKey, err := model.NewDataKey()
		if err != nil {
			return nil, errors.Wrap(err, op)
		}

		key, err = k.keyRepo.RotateKey(ctx, newKey)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}
	}

	sealed, err := k.cipher.Seal(secret, key.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	sealed.KeyID = key.ID

	dataSize := int64(len(secret.Data))
	err = k.keyRepo.UpdateStats(ctx, key.ID, dataSize)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return sealed, nil
}

func (k *keyService) Unseal(ctx context.Context, secret *model.Secret) (*model.Secret, error) {
	const op = "unseal"

	dataKey, err := k.keyRepo.GetByID(ctx, secret.KeyID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	unsealed, err := k.cipher.Unseal(secret, dataKey.Key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}
