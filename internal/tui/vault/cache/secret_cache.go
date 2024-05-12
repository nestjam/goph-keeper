package cache

import (
	vault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type SecretsCache struct {
	secrets map[string]*vault.Secret
}

func New() *SecretsCache {
	return &SecretsCache{
		secrets: make(map[string]*vault.Secret),
	}
}

func (c *SecretsCache) CacheSecrets(secrets []*vault.Secret) {
	for i := 0; i < len(secrets); i++ {
		secret := secrets[i]
		if _, ok := c.secrets[secret.ID]; !ok {
			c.secrets[secret.ID] = secret
		}
	}
}

func (c *SecretsCache) ListSecrets() []*vault.Secret {
	secrets := make([]*vault.Secret, len(c.secrets))

	i := 0
	for _, secret := range c.secrets {
		secrets[i] = secret
		i++
	}

	return secrets
}

func (c *SecretsCache) CacheSecret(secret *vault.Secret) {
	c.secrets[secret.ID] = secret
}

func (c *SecretsCache) GetSecret(id string) (*vault.Secret, bool) {
	secret, ok := c.secrets[id]
	return secret, ok
}

func (c *SecretsCache) RemoveSecret(id string) error {
	panic("unimplemented")
}
