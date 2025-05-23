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
	cipher  *model.MasterKeyCipher
	config  KeyRotationConfig
}

func NewKeyService(keyRepo vault.DataKeyRepository, config KeyRotationConfig, rootKey *model.MasterKey) *keyService {
	return &keyService{
		keyRepo: keyRepo,
		config:  config,
		cipher:  model.NewMasterKeyCipher(rootKey),
	}
}

func (k *keyService) Seal(ctx context.Context, secret *model.Secret) (*model.Secret, error) {
	const op = "seal"

	key, err := k.keyRepo.GetKey(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	if key == nil ||
		key.EncryptedDataSize >= k.config.EncryptedDataSizeThreshold ||
		key.EncryptionsCount >= k.config.EncryptionsCountThreshold {
		key, err = k.rotateKey(ctx)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}
	}

	key, err = k.cipher.Unseal(key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	cipher := model.NewDataKeyCipher(key)
	sealed, err := cipher.Seal(secret)
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

func (k *keyService) rotateKey(ctx context.Context) (*model.DataKey, error) {
	const op = "rotate data key"

	key, err := model.NewDataKey()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	key, err = k.cipher.Seal(key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	key, err = k.keyRepo.RotateKey(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return key, nil
}

func (k *keyService) Unseal(ctx context.Context, secret *model.Secret) (*model.Secret, error) {
	const op = "unseal"

	key, err := k.keyRepo.GetByID(ctx, secret.KeyID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	key, err = k.cipher.Unseal(key)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	cipher := model.NewDataKeyCipher(key)
	unsealed, err := cipher.Unseal(secret)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return unsealed, nil
}
