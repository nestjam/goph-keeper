package cache

import (
	vault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type secretCache struct {
	*vault.Secret
	dataCached bool
}

type SecretsCache struct {
	secrets map[string]secretCache
}

func New() *SecretsCache {
	return &SecretsCache{
		secrets: make(map[string]secretCache),
	}
}

func (c *SecretsCache) CacheSecrets(secrets []*vault.Secret) {
	newCache := make(map[string]secretCache, len(c.secrets))

	for i := 0; i < len(secrets); i++ {
		secret := secrets[i]
		if cached, ok := c.secrets[secret.ID]; ok {
			newCache[secret.ID] = cached
			continue
		}
		newCache[secret.ID] = secretCache{Secret: secret}
	}

	c.secrets = newCache
}

func (c *SecretsCache) ListSecrets() []*vault.Secret {
	secrets := make([]*vault.Secret, len(c.secrets))

	i := 0
	for _, cache := range c.secrets {
		secrets[i] = cache.Secret
		i++
	}

	return secrets
}

func (c *SecretsCache) CacheSecret(secret *vault.Secret) {
	cache := secretCache{
		Secret:     secret,
		dataCached: true,
	}
	c.secrets[secret.ID] = cache
}

func (c *SecretsCache) GetSecret(id string) (secret *vault.Secret, dataCached bool, found bool) {
	if cache, ok := c.secrets[id]; ok {
		found = true
		secret = cache.Secret
		dataCached = cache.dataCached
	}
	return
}

func (c *SecretsCache) RemoveSecret(id string) {
	delete(c.secrets, id)
}
