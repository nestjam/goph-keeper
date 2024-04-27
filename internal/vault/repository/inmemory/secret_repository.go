package inmemory

import (
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

func (r *secretRepository) ListSecrets(userID uuid.UUID) ([]*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	secrets := r.userSecrets[userID]
	return secrets, nil
}

func (r *secretRepository) AddSecret(secret *model.Secret, userID uuid.UUID) (*model.Secret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New()
	createdSecret := &model.Secret{ID: id}
	r.secrets[createdSecret.ID] = createdSecret

	addUserSecret(r, createdSecret, userID)

	return createdSecret, nil
}

func addUserSecret(r *secretRepository, secret *model.Secret, userID uuid.UUID) {
	secrets := r.userSecrets[userID]
	secrets = append(secrets, secret)
	r.userSecrets[userID] = secrets
}
