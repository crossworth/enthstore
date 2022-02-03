## Hstore type and predicates for Ent.

This package defines the [`Postgres hstore`](https://www.postgresql.org/docs/9/hstore.html) type
to be used with [`ent`](https://github.com/ent/ent) and a few predicates.

### Using the type:
```go
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Other("attributes", enthstore.Hstore{}).
			SchemaType(enthstore.Hstore{}.SchemaType()).
			Default(func() enthstore.Hstore {
				return enthstore.Hstore{}
			}),
	}
}
```

### Using the predicates:
```go
users, err := client.User.Query().Where(func(selector *sql.Selector) {
    selector.Where(enthstore.HasKey(user.FieldAttributes, "a"))
}).All(context.Background)
```

#### List of predicates:
- HasKey
- HasAllKeys
- ValueIsNull
- ValueEQ
- ValueNEQ
- ValueGT
- ValueGTE
- ValueLT
- ValueLTE
- ValueContains
- ValueHasPrefix
- ValueHasSuffix

### Using with [GQLGen](https://github.com/99designs/gqlgen):

Define a [custom scalar](https://gqlgen.com/reference/scalars/):
```graphql
scalar Hstore
```

And declare the type mapping on `gqlgen.yml`:

```yaml
models:
  Hstore:
    model: github.com/crossworth/enthstore.Hstore
```