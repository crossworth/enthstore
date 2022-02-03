package enthstore

import (
	"strings"

	"entgo.io/ent/dialect/sql"
)

// HasKey checks if the given column has the provided key.
func HasKey(column string, key string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.WriteString("exist(").Ident(column).Comma().WriteString(quoteKey(key)).WriteString(")")
	})
}

// HasAllKeys checks if the given column has all the keys provided.
func HasAllKeys(column string, keys ...string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(column).WriteString(" ?& ").
			WriteString("ARRAY[")

		quoted := make([]string, 0, len(keys))
		for _, k := range keys {
			quoted = append(quoted, quoteKey(k))
		}

		b.WriteString(strings.Join(quoted, ",")).WriteString("]")
	})
}

// ValueIsNull check if the given column has a key which the value is null.
func ValueIsNull(column string, key string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.WriteString("defined(").Ident(column).Comma().WriteString(quoteKey(key)).WriteString(") is false")
	})
}

// ValueEQ check if the given column has a key which the value is equals to the provided string.
func ValueEQ(column string, key string, val string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).WriteOp(sql.OpEQ).Arg(val)
	})
}

// ValueNEQ check if the given column has a key which the value is not equals to the provided string.
func ValueNEQ(column string, key string, val string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).WriteOp(sql.OpNEQ).Arg(val)
	})
}

// ValueGT check if the given column has a key which the value is greater than the provided string.
func ValueGT(column string, key string, val string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).WriteOp(sql.OpGT).Arg(val)
	})
}

// ValueGTE check if the given column has a key which the value is greater
// or equals to the provided string.
func ValueGTE(column string, key string, val string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).WriteOp(sql.OpGTE).Arg(val)
	})
}

// ValueLT check if the given column has a key which the value is smaller than the provided string.
func ValueLT(column string, key string, val string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).WriteOp(sql.OpLT).Arg(val)
	})
}

// ValueLTE check if the given column has a key which the value is smaller
// or equals to the provided string.
func ValueLTE(column string, key string, val string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).WriteOp(sql.OpLTE).Arg(val)
	})
}

// ValueContains check given column has a key which the value contains the provided string.
func ValueContains(column string, key, val string) *sql.Predicate {
	return sql.P().Contains(sql.P().Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).String(), val)
}

// ValueHasPrefix check given column has a key which the value the provided prefix.
func ValueHasPrefix(column string, key, val string) *sql.Predicate {
	return sql.P().HasPrefix(sql.P().Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).String(), val)
}

// ValueHasSuffix check given column has a key which the value the provided suffix.
func ValueHasSuffix(column string, key, val string) *sql.Predicate {
	return sql.P().HasSuffix(sql.P().Ident(column).WriteString(" -> ").WriteString(quoteKey(key)).String(), val)
}
