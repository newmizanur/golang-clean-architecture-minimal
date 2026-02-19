package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Contact struct {
	ent.Schema
}

func (Contact) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").MaxLen(36).NotEmpty(),
		field.String("first_name").MaxLen(100).NotEmpty(),
		field.String("last_name").MaxLen(100).Optional().Nillable(),
		field.String("email").MaxLen(100).Optional().Nillable(),
		field.String("phone").MaxLen(100).Optional().Nillable(),
		field.String("user_id").MaxLen(36).NotEmpty(),
		field.Int64("created_at"),
		field.Int64("updated_at"),
	}
}

func (Contact) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("contacts").
			Field("user_id").
			Required().
			Unique(),
		edge.To("addresses", Address.Type),
	}
}
