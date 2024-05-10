//go:build integration

package secret

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	modelAuth "github.com/nestjam/goph-keeper/internal/auth/model"
	"github.com/nestjam/goph-keeper/internal/auth/repository/pgsql"
	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
	modelVault "github.com/nestjam/goph-keeper/internal/vault/model"
	"github.com/nestjam/goph-keeper/internal/vault/repository/pgsql/key"
	"github.com/nestjam/goph-keeper/migration"
)

var h *utils.PGSQLRepositoryTestHelper

func TestMain(m *testing.M) {
	h = &utils.PGSQLRepositoryTestHelper{}
	h.Run(m)
}

func TestSecretRepository(t *testing.T) {
	vault.SecretRepositoryContract{
		NewSecretRepository: func() (vault.SecretRepository, func(), vault.SecretTestData) {
			t.Helper()

			ctx := context.Background()
			var err error
			r, err := NewSecretRepository(ctx, h.DataSourceName)
			require.NoError(t, err)

			closer := func() {
				r.Close()

				migrator := migration.NewDatabaseMigrator(h.DataSourceName)
				_ = migrator.Drop()
			}

			testData := vault.SecretTestData{
				Users: setupUsers(t),
				Keys:  setupKeys(t),
			}
			return r, closer, testData
		},
	}.Test(t)
}

func setupUsers(t *testing.T) uuid.UUIDs {
	t.Helper()

	ctx := context.Background()
	r, err := pgsql.NewUserRepository(ctx, h.DataSourceName)
	require.NoError(t, err)

	user, err := r.Register(ctx, modelAuth.User{Email: "user@email.com", Password: "1"})
	require.NoError(t, err)
	user2, err := r.Register(ctx, modelAuth.User{Email: "user2@email.com", Password: "2"})
	require.NoError(t, err)

	return uuid.UUIDs{user.ID, user2.ID}
}

func setupKeys(t *testing.T) uuid.UUIDs {
	t.Helper()

	ctx := context.Background()
	r, err := key.NewDataKeyRepository(ctx, h.DataSourceName)
	require.NoError(t, err)

	key, err := r.RotateKey(ctx, &modelVault.DataKey{})
	require.NoError(t, err)

	return uuid.UUIDs{key.ID}
}
