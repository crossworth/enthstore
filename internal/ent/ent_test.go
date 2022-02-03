package ent

import (
	"context"
	"testing"

	"internal/databasetest"
	"internal/ent/ent"
	"internal/ent/ent/user"

	"entgo.io/ent/dialect/sql"
	"github.com/crossworth/enthstore"
	"github.com/stretchr/testify/require"
)

func TestIntegrationHstorePGX(t *testing.T) {
	databasetest.RunWithEnt(t, "pgx", func(client *ent.Client) {

		t.Run("empty", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			u := client.User.Create().SaveX(context.Background())

			require.Len(t, u.Attributes, 0)
		})

		t.Run("simple values", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			u := client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())

			require.Len(t, u.Attributes, 2)
		})

		t.Run("HasKey (no key)", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"d": "b",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.HasKey(user.FieldAttributes, "a"))
			}).CountX(context.Background())
			require.Equal(t, 0, count)
		})

		t.Run("HasKey", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.HasKey(user.FieldAttributes, "a"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("HasAllKeys", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.HasAllKeys(user.FieldAttributes, "a", "c"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueIsNull", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.Hstore{
				"a": nil,
			}).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueIsNull(user.FieldAttributes, "a"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueEQ", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.Hstore{
				"a": nil,
			}).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueEQ(user.FieldAttributes, "a", "b"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueNEQ", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "1",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueNEQ(user.FieldAttributes, "a", "b"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueGT", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a9",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueGT(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueGTE", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueGTE(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueLT", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a9",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueLT(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueLTE", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueLTE(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 2, count)
		})

		t.Run("ValueContains", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "g",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueContains(user.FieldAttributes, "a", "a1"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueHasPrefix", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "2a",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueHasPrefix(user.FieldAttributes, "a", "a"))
			}).CountX(context.Background())
			require.Equal(t, 2, count)
		})

		t.Run("ValueHasSuffix", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueHasSuffix(user.FieldAttributes, "a", "d"))
			}).CountX(context.Background())
			require.Equal(t, 2, count)
		})
	})
}

func TestIntegrationHstorePQ(t *testing.T) {
	databasetest.RunWithEnt(t, "postgres", func(client *ent.Client) {

		t.Run("empty", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			u := client.User.Create().SaveX(context.Background())

			require.Len(t, u.Attributes, 0)
		})

		t.Run("simple values", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			u := client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())

			require.Len(t, u.Attributes, 2)
		})

		t.Run("HasKey (no key)", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"d": "b",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.HasKey(user.FieldAttributes, "a"))
			}).CountX(context.Background())
			require.Equal(t, 0, count)
		})

		t.Run("HasKey", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.HasKey(user.FieldAttributes, "a"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("HasAllKeys", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.HasAllKeys(user.FieldAttributes, "a", "c"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueIsNull", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.Hstore{
				"a": nil,
			}).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueIsNull(user.FieldAttributes, "a"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueEQ", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.Hstore{
				"a": nil,
			}).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueEQ(user.FieldAttributes, "a", "b"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueNEQ", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "b",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "1",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueNEQ(user.FieldAttributes, "a", "b"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueGT", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a9",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueGT(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueGTE", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueGTE(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueLT", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a9",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueLT(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueLTE", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueLTE(user.FieldAttributes, "a", "a2"))
			}).CountX(context.Background())
			require.Equal(t, 2, count)
		})

		t.Run("ValueContains", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "g",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueContains(user.FieldAttributes, "a", "a1"))
			}).CountX(context.Background())
			require.Equal(t, 1, count)
		})

		t.Run("ValueHasPrefix", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
				"c": "d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "2a",
				"c": "d",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueHasPrefix(user.FieldAttributes, "a", "a"))
			}).CountX(context.Background())
			require.Equal(t, 2, count)
		})

		t.Run("ValueHasSuffix", func(t *testing.T) {
			defer client.User.Delete().ExecX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a1d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2d",
			})).SaveX(context.Background())
			client.User.Create().SetAttributes(enthstore.FromMap(map[string]string{
				"a": "a2",
			})).SaveX(context.Background())

			count := client.User.Query().Where(func(selector *sql.Selector) {
				selector.Where(enthstore.ValueHasSuffix(user.FieldAttributes, "a", "d"))
			}).CountX(context.Background())
			require.Equal(t, 2, count)
		})
	})
}
