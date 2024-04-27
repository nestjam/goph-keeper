package inmemory

import (
	"testing"

	"github.com/nestjam/goph-keeper/internal/vault"
)

func TestSecretRepository(t *testing.T) {
	vault.SecretRepositoryContract{
		NewSecretRepository: func() (vault.SecretRepository, func()) {
			t.Helper()

			r := NewSecretRepository()
			return r, func() {
			}
		},
	}.Test(t)
}
