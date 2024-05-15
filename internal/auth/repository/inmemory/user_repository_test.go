package inmemory

import (
	"testing"

	"github.com/nestjam/goph-keeper/internal/auth"
)

func TestUserRepository(t *testing.T) {
	auth.UserRepositoryContract{
		NewUserRepository: func() (auth.UserRepository, func()) {
			t.Helper()

			r := NewUserRepository()
			return r, func() {
			}
		},
	}.Test(t)
}
