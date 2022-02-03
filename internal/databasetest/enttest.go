package databasetest

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"internal/ent/ent"
	"internal/ent/ent/enttest"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
)

// RunWithEnt runs a test with a database and *ent.Client.
func RunWithEnt(t *testing.T, driver string, fn func(client *ent.Client)) {
	t.Helper()
	RunWithDatabase(t, driver, func(db *sql.DB, purgeDB func()) {
		purgeDB()

		drv := entsql.OpenDB(dialect.Postgres, db)
		client := enttest.NewClient(t,
			enttest.WithOptions(
				ent.Driver(drv),
			),
			enttest.WithMigrateOptions(
				EnableHstoreOption(db),
			),
		)

		defer func() {
			_ = client.Close()
		}()

		// runs the test
		fn(client)
	})
}

// EnableHstoreOption returns a schema.MigrateOption
// that will enable the Postgres hstore extension if needed.
func EnableHstoreOption(db *sql.DB) schema.MigrateOption {
	return schema.WithHooks(func(next schema.Creator) schema.Creator {
		return schema.CreateFunc(func(ctx context.Context, tables ...*schema.Table) error {
			_, err := db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS hstore WITH SCHEMA public;`)
			if err != nil {
				return fmt.Errorf("could not enable hstore extension: %w", err)
			}

			return next.Create(ctx, tables...)
		})
	})
}
