package schema

import (
	"time"

	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ProjectUser struct {
	ent.Schema
}

func (ProjectUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.Enum("role").
			Values("owner", "member", "viewer").
			Default("member"),

		field.Time("created_at").Default(time.Now),
	}
}

func (ProjectUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("memberships").
			Unique().
			Required(),

		edge.From("user", User.Type).
			Ref("memberships").
			Unique().
			Required(),
	}
}

func (ProjectUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("project", "user").Unique(),
	}
}
