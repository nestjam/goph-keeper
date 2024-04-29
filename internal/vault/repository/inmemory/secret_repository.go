package inmemory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type userSecrets map[uuid.UUID]*model.Secret

type secretRepository struct {
	userSecrets map[uuid.UUID]userSecrets
	mu          sync.Mutex
}

func NewSecretRepository() vault.SecretRepository {
	return &secretRepository{
		userSecrets: make(map[uuid.UUID]userSecrets),
	}
}

func (r *secretRepository) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userSecrets := r.userSecrets[userID]
	secrets := make([]*model.Secret, len(userSecrets))

	i := 0
	for _, secret := range userSecrets {
		secrets[i] = copySecret(secret)
		secrets[i].Data = ""
		i++
	}

	return secrets, nil
}

func (r *secretRepository) AddSecret(ctx context.Context, s *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	secret := copySecret(s)
	secret.ID = uuid.New()

	if _, ok := r.userSecrets[userID]; !ok {
		r.userSecrets[userID] = make(userSecrets)
	}
	secrets := r.userSecrets[userID]
	secrets[secret.ID] = secret

	return secret, nil
}

func (r *secretRepository) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userSecrets, ok := r.userSecrets[userID]
	if !ok {
		return nil, vault.ErrUserDoesNotExist
	}

	secret, ok := userSecrets[secretID]
	if !ok {
		return nil, vault.ErrSecretDoesNotExist
	}

	return secret, nil
}

func (r *secretRepository) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userSecrets, ok := r.userSecrets[userID]
	if !ok {
		return vault.ErrUserDoesNotExist
	}

	delete(userSecrets, secretID)

	return nil
}

func copySecret(secret *model.Secret) *model.Secret {
	return &model.Secret{
		ID:   secret.ID,
		Data: secret.Data,
	}
}
