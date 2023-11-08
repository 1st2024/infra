package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Team struct {
	ent.Schema
}

func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Immutable().Unique(),
		field.Time("created_at").Immutable().Default(time.Now).Annotations(
			entsql.Default("CURRENT_TIMESTAMP"),
		),
		field.Bool("is_default"),
		field.Bool("is_blocked"),
		field.String("name"),
		field.String("tier"),
	}
}

func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type).Through("users_teams", UsersTeams.Type),
		edge.To("team_api_keys", TeamApiKey.Type),
		edge.From("team_tier", Tier.Type).Ref("teams").Unique().Field("tier").Required(),
		edge.To("envs", Env.Type),
	}
}
