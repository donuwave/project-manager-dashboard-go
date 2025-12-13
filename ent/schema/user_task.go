package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserTask struct {
	ent.Schema
}

func (UserTask) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.Time("created_at").Default(time.Now),
	}
}

func (UserTask) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("assignments").
			Unique().
			Required(),

		edge.From("task", Task.Type).
			Ref("assignments").
			Unique().
			Required(),
	}
}

func (UserTask) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user", "task").Unique(),
	}
}
