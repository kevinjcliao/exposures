// Code generated by entc, DO NOT EDIT.

package checkin

import (
	"helloworld/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CheckinTime applies equality check predicate on the "checkin_time" field. It's identical to CheckinTimeEQ.
func CheckinTime(v int64) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCheckinTime), v))
	})
}

// EventID applies equality check predicate on the "event_id" field. It's identical to EventIDEQ.
func EventID(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEventID), v))
	})
}

// CheckinTimeEQ applies the EQ predicate on the "checkin_time" field.
func CheckinTimeEQ(v int64) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCheckinTime), v))
	})
}

// CheckinTimeNEQ applies the NEQ predicate on the "checkin_time" field.
func CheckinTimeNEQ(v int64) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCheckinTime), v))
	})
}

// CheckinTimeIn applies the In predicate on the "checkin_time" field.
func CheckinTimeIn(vs ...int64) predicate.Checkin {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Checkin(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCheckinTime), v...))
	})
}

// CheckinTimeNotIn applies the NotIn predicate on the "checkin_time" field.
func CheckinTimeNotIn(vs ...int64) predicate.Checkin {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Checkin(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCheckinTime), v...))
	})
}

// CheckinTimeGT applies the GT predicate on the "checkin_time" field.
func CheckinTimeGT(v int64) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCheckinTime), v))
	})
}

// CheckinTimeGTE applies the GTE predicate on the "checkin_time" field.
func CheckinTimeGTE(v int64) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCheckinTime), v))
	})
}

// CheckinTimeLT applies the LT predicate on the "checkin_time" field.
func CheckinTimeLT(v int64) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCheckinTime), v))
	})
}

// CheckinTimeLTE applies the LTE predicate on the "checkin_time" field.
func CheckinTimeLTE(v int64) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCheckinTime), v))
	})
}

// EventIDEQ applies the EQ predicate on the "event_id" field.
func EventIDEQ(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEventID), v))
	})
}

// EventIDNEQ applies the NEQ predicate on the "event_id" field.
func EventIDNEQ(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEventID), v))
	})
}

// EventIDIn applies the In predicate on the "event_id" field.
func EventIDIn(vs ...string) predicate.Checkin {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Checkin(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldEventID), v...))
	})
}

// EventIDNotIn applies the NotIn predicate on the "event_id" field.
func EventIDNotIn(vs ...string) predicate.Checkin {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Checkin(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldEventID), v...))
	})
}

// EventIDGT applies the GT predicate on the "event_id" field.
func EventIDGT(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldEventID), v))
	})
}

// EventIDGTE applies the GTE predicate on the "event_id" field.
func EventIDGTE(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldEventID), v))
	})
}

// EventIDLT applies the LT predicate on the "event_id" field.
func EventIDLT(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldEventID), v))
	})
}

// EventIDLTE applies the LTE predicate on the "event_id" field.
func EventIDLTE(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldEventID), v))
	})
}

// EventIDContains applies the Contains predicate on the "event_id" field.
func EventIDContains(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldEventID), v))
	})
}

// EventIDHasPrefix applies the HasPrefix predicate on the "event_id" field.
func EventIDHasPrefix(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldEventID), v))
	})
}

// EventIDHasSuffix applies the HasSuffix predicate on the "event_id" field.
func EventIDHasSuffix(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldEventID), v))
	})
}

// EventIDEqualFold applies the EqualFold predicate on the "event_id" field.
func EventIDEqualFold(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldEventID), v))
	})
}

// EventIDContainsFold applies the ContainsFold predicate on the "event_id" field.
func EventIDContainsFold(v string) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldEventID), v))
	})
}

// HasSender applies the HasEdge predicate on the "sender" edge.
func HasSender() predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SenderTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, SenderTable, SenderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSenderWith applies the HasEdge predicate on the "sender" edge with a given conditions (other predicates).
func HasSenderWith(preds ...predicate.User) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SenderInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, SenderTable, SenderColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Checkin) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Checkin) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Checkin) predicate.Checkin {
	return predicate.Checkin(func(s *sql.Selector) {
		p(s.Not())
	})
}
