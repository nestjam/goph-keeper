package inmemory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nestjam/goph-keeper/internal/vault"
)

func TestSecretRepository(t *testing.T) {
	vault.SecretRepositoryContract{
		NewSecretRepository: func() (vault.SecretRepository, func(), vault.SecretTestData) {
			t.Helper()

			r := NewSecretRepository()
			closer := func() {}
			testData := vault.SecretTestData{
				Users: uuid.UUIDs{uuid.New(), uuid.New()},
				Keys:  uuid.UUIDs{uuid.New()},
			}
			return r, closer, testData
		},
	}.Test(t)
}
