package inmemory

import (
	"testing"

	"github.com/nestjam/goph-keeper/internal/vault"
)

func TestDataKeyRepository(t *testing.T) {
	vault.DataKeyRepositoryContract{
		NewDataKeyRepository: func() (vault.DataKeyRepository, func()) {
			t.Helper()

			r := NewDataKeyRepository()
			return r, func() {
			}
		},
	}.Test(t)
}
