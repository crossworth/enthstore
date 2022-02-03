package enthstore

import (
	"strconv"
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/require"
)

func TestHstorePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input     sql.Querier
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(HasKey("attributes", "'test'")),
			wantQuery: `SELECT * FROM "users" WHERE exist("attributes", '''test''')`,
			wantArgs:  nil,
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(HasKey("attributes", "test")),
			wantQuery: `SELECT * FROM "users" WHERE exist("attributes", 'test')`,
			wantArgs:  nil,
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(HasAllKeys("attributes", "test", "test1")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" ?& ARRAY['test','test1']`,
			wantArgs:  nil,
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueIsNull("attributes", "test")),
			wantQuery: `SELECT * FROM "users" WHERE defined("attributes", 'test') is false`,
			wantArgs:  nil,
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueEQ("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' = $1`,
			wantArgs:  []interface{}{"val"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueNEQ("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' <> $1`,
			wantArgs:  []interface{}{"val"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueGT("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' > $1`,
			wantArgs:  []interface{}{"val"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueGTE("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' >= $1`,
			wantArgs:  []interface{}{"val"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueLT("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' < $1`,
			wantArgs:  []interface{}{"val"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueLTE("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' <= $1`,
			wantArgs:  []interface{}{"val"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueContains("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' LIKE $1`,
			wantArgs:  []interface{}{"%val%"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueHasPrefix("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' LIKE $1`,
			wantArgs:  []interface{}{"val%"},
		},
		{
			input: sql.Dialect(dialect.Postgres).
				Select("*").
				From(sql.Table("users")).
				Where(ValueHasSuffix("attributes", "key", "val")),
			wantQuery: `SELECT * FROM "users" WHERE "attributes" -> 'key' LIKE $1`,
			wantArgs:  []interface{}{"%val"},
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			query, args := tt.input.Query()
			require.Equal(t, tt.wantQuery, query)
			require.Equal(t, tt.wantArgs, args)
		})
	}
}
