package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ProjectTask struct {
	ent.Schema
}

func (ProjectTask) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.Time("created_at").Default(time.Now),
	}
}

func (ProjectTask) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("project_tasks").
			Unique().
			Required(),

		edge.From("task", Task.Type).
			Ref("project_tasks").
			Unique().
			Required(),
	}
}

func (ProjectTask) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("project", "task").Unique(),
	}
}
