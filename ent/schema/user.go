package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.String("email").Unique(),
		field.String("name"),
		field.String("country").Optional(),

		field.Time("created_at").Default(time.Now),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assigned_tasks", Task.Type),
		edge.To("memberships", ProjectUser.Type),
	}
}
