package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Item struct {
	ent.Schema
}

func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.String("sku").MaxLen(255).NotEmpty().Unique(),
		field.String("currency").MaxLen(10).NotEmpty(),
		field.Int32("stock"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Item) Edges() []ent.Edge {
	return nil
}
