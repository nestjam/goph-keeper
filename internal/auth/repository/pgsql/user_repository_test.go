//go:build integration

package pgsql

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/utils"
	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/migration"
)

var h *utils.PGSQLRepositoryTestHelper

func TestMain(m *testing.M) {
	h = &utils.PGSQLRepositoryTestHelper{}
	h.Run(m)
}

func TestUserRepository(t *testing.T) {
	auth.UserRepositoryContract{
		NewUserRepository: func() (auth.UserRepository, func()) {
			t.Helper()

			ctx := context.Background()
			r, err := NewUserRepository(ctx, h.DataSourceName)
			require.NoError(t, err)

			return r, func() {
				r.Close()

				migrator := migration.NewDatabaseMigrator(h.DataSourceName)
				_ = migrator.Drop()
			}
		},
	}.Test(t)
}
