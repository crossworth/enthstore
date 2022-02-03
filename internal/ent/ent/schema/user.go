package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/crossworth/enthstore"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Other("attributes", enthstore.Hstore{}).
			SchemaType(enthstore.Hstore{}.SchemaType()).
			Default(func() enthstore.Hstore {
				return enthstore.Hstore{}
			}),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
