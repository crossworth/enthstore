package hstore

import (
	"database/sql"
	"testing"

	"internal/databasetest"

	"github.com/crossworth/enthstore"
	"github.com/stretchr/testify/require"
)

func TestIntegrationHstorePGX(t *testing.T) {
	databasetest.RunWithDatabase(t, "pgx", func(db *sql.DB, purgeDB func()) {
		_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS hstore WITH SCHEMA public;")
		require.NoError(t, err)

		hs := enthstore.Hstore{}

		t.Run("null", func(t *testing.T) {
			err = db.QueryRow("SELECT NULL::hstore").Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("null value", func(t *testing.T) {
			err = db.QueryRow("SELECT $1::hstore", enthstore.Hstore{}).Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("null ptr value", func(t *testing.T) {
			err = db.QueryRow("SELECT $1::hstore", &enthstore.Hstore{}).Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("empty hstore", func(t *testing.T) {
			err = db.QueryRow("SELECT ''::hstore").Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("key value", func(t *testing.T) {
			err = db.QueryRow("SELECT ''::hstore").Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("all kind of values", func(t *testing.T) {
			input := enthstore.FromMap(map[string]string{
				"k1":       "v1",
				"NOT NULL": "NULL",
				"a":        "a",
				"a'a":      `b'b`,
				`"a"`:      `"b"`,
				"tes t":    `test test`,
				"hs":       "a=>b",
				"hs2":      `"a"=>b"`,
				"a\tb":     "\nabc\t",
				"empty":    "",
			})
			input.Set("NULL", nil)

			output := enthstore.Hstore{}

			err = db.QueryRow("SELECT $1::hstore", input).Scan(&output)
			require.NoError(t, err)
			require.True(t, input.Equals(output))
		})
	})
}

func TestIntegrationHstorePQ(t *testing.T) {
	databasetest.RunWithDatabase(t, "postgres", func(db *sql.DB, purgeDB func()) {
		_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS hstore WITH SCHEMA public;")
		require.NoError(t, err)

		hs := enthstore.Hstore{}

		t.Run("null", func(t *testing.T) {
			err = db.QueryRow("SELECT NULL::hstore").Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("null value", func(t *testing.T) {
			err = db.QueryRow("SELECT $1::hstore", enthstore.Hstore{}).Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("null ptr value", func(t *testing.T) {
			err = db.QueryRow("SELECT $1::hstore", &enthstore.Hstore{}).Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("empty hstore", func(t *testing.T) {
			err = db.QueryRow("SELECT ''::hstore").Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("key value", func(t *testing.T) {
			err = db.QueryRow("SELECT ''::hstore").Scan(&hs)
			require.NoError(t, err)
			require.Len(t, hs, 0)
		})

		t.Run("all kind of values", func(t *testing.T) {
			input := enthstore.FromMap(map[string]string{
				"k1":       "v1",
				"NOT NULL": "NULL",
				"a":        "a",
				"a'a":      `b'b`,
				`"a"`:      `"b"`,
				"tes t":    `test test`,
				"hs":       "a=>b",
				"hs2":      `"a"=>b"`,
				"a\tb":     "\nabc\t",
				"empty":    "",
			})
			input.Set("NULL", nil)

			output := enthstore.Hstore{}

			err = db.QueryRow("SELECT $1::hstore", input).Scan(&output)
			require.NoError(t, err)
			require.True(t, input.Equals(output))
		})
	})
}
