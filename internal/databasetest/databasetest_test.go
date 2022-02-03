package databasetest

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationRunWithDatabase(t *testing.T) {
	RunWithDatabase(t, "pgx", func(db *sql.DB, purgeDB func()) {
		row := db.QueryRow(`SELECT 1`)
		require.Nil(t, row.Err())

		var result int
		err := row.Scan(&result)
		require.NoError(t, err)
		require.Equal(t, 1, result)
	})
}
