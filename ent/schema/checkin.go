package schema

import (
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
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
			Validate(func(eventName string) error {
				if _, err := uuid.Parse(eventName); err != nil {
					return errors.New("event name is not well-formed UUID")
				}
				return nil
			}),
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
