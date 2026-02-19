package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Address struct {
	ent.Schema
}

func (Address) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").MaxLen(36).NotEmpty(),
		field.String("contact_id").MaxLen(36).NotEmpty(),
		field.String("street").MaxLen(255).Optional().Nillable(),
		field.String("city").MaxLen(255).Optional().Nillable(),
		field.String("province").MaxLen(255).Optional().Nillable(),
		field.String("postal_code").MaxLen(10).Optional().Nillable(),
		field.String("country").MaxLen(100).Optional().Nillable(),
		field.Int64("created_at"),
		field.Int64("updated_at"),
	}
}

func (Address) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("contact", Contact.Type).
			Ref("addresses").
			Field("contact_id").
			Required().
			Unique(),
	}
}
