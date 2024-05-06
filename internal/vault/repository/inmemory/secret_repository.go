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
		secrets[i] = secret.Copy()
		secrets[i].Data = nil
		i++
	}

	return secrets, nil
}

func (r *secretRepository) AddSecret(ctx context.Context, s *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	secret := s.Copy()
	secret.ID = uuid.New()
	return r.addOrUpdateSecret(secret, userID), nil
}

func (r *secretRepository) UpdateSecret(ctx context.Context, s *model.Secret, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	secret := s.Copy()
	_ = r.addOrUpdateSecret(secret, userID)
	return nil
}

func (r *secretRepository) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userSecrets, ok := r.userSecrets[userID]
	if !ok {
		return nil, vault.ErrSecretNotFound
	}

	secret, ok := userSecrets[secretID]
	if !ok {
		return nil, vault.ErrSecretNotFound
	}

	return secret, nil
}

func (r *secretRepository) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userSecrets, ok := r.userSecrets[userID]
	if !ok {
		return nil
	}

	delete(userSecrets, secretID)

	return nil
}

func (r *secretRepository) addOrUpdateSecret(secret *model.Secret, userID uuid.UUID) *model.Secret {
	if _, ok := r.userSecrets[userID]; !ok {
		r.userSecrets[userID] = make(userSecrets)
	}
	secrets := r.userSecrets[userID]
	secrets[secret.ID] = secret

	return secret
}
