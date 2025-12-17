package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Task struct {
	ent.Schema
}

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.String("title"),
		field.String("description").Optional(),

		field.Enum("status").
			Values("todo", "in_progress", "done").
			Default("todo"),

		field.Enum("priority").
			Values("low", "medium", "high").
			Default("medium"),

		field.Int("position").Default(0),

		field.Time("due_date").Optional(),

		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.UUID("assignee_id", uuid.UUID{}).Optional().Nillable(),
	}
}

func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project_tasks", ProjectTask.Type),
		edge.From("assignee", User.Type).
			Ref("assigned_tasks").
			Field("assignee_id").
			Unique(),
	}
}
