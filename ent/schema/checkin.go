package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Checkin holds the schema definition for the Checkin entity.
type Checkin struct {
	ent.Schema
}

// Fields of the Checkin.
func (Checkin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("checkin_time").
			Positive(),
		field.String("event_id").
			MinLen(10),
	}
}

// Edges of the Checkin.
func (Checkin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sender", User.Type).
			Ref("checkins").
			Unique(),
	}
}
