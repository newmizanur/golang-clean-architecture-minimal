package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").MaxLen(36).NotEmpty(),
		field.String("name").MaxLen(100).NotEmpty(),
		field.String("password").MaxLen(100).NotEmpty(),
		field.Int64("created_at"),
		field.Int64("updated_at"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contacts", Contact.Type),
	}
}
