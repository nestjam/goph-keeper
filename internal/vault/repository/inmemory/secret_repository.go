package inmemory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/internal/vault/model"
)

type secretRepository struct {
	secrets     map[uuid.UUID]*model.Secret
	userSecrets map[uuid.UUID][]*model.Secret
	mu          sync.Mutex
}

func NewSecretRepository() vault.SecretRepository {
	return &secretRepository{
		secrets:     make(map[uuid.UUID]*model.Secret),
		userSecrets: make(map[uuid.UUID][]*model.Secret),
	}
}

func (r *secretRepository) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	secrets := r.userSecrets[userID]
	return secrets, nil
}

func (r *secretRepository) AddSecret(ctx context.Context, s *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	secret := copySecret(s)
	secret.ID = uuid.New()

	r.secrets[secret.ID] = secret
	addUserSecret(r, secret, userID)

	return secret, nil
}

func copySecret(secret *model.Secret) *model.Secret {
	return &model.Secret{
		Payload: secret.Payload,
	}
}

func addUserSecret(r *secretRepository, secret *model.Secret, userID uuid.UUID) {
	secrets := r.userSecrets[userID]
	secrets = append(secrets, secret)
	r.userSecrets[userID] = secrets
}
