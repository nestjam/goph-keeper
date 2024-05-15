//go:build integration

package key

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/vault"
	"github.com/nestjam/goph-keeper/migration"
)

var h *utils.PGSQLRepositoryTestHelper

func TestMain(m *testing.M) {
	h = &utils.PGSQLRepositoryTestHelper{}
	h.Run(m)
}

func TestKeyRepository(t *testing.T) {
	vault.DataKeyRepositoryContract{
		NewDataKeyRepository: func() (vault.DataKeyRepository, func()) {
			t.Helper()

			dsn := h.DataSourceName
			migrator := migration.NewDatabaseMigrator(dsn)
			err := migrator.Up()
			require.NoError(t, err)

			ctx := context.Background()
			r, err := NewDataKeyRepository(ctx, dsn)
			require.NoError(t, err)

			return r, func() {
				r.Close()

				migrator := migration.NewDatabaseMigrator(dsn)
				_ = migrator.Drop()
			}
		},
	}.Test(t)
}
