package schema

import "entgo.io/ent"

// OperatorProfile holds the schema definition for the OperatorProfile entity.
type OperatorProfile struct {
	ent.Schema
}

// Fields of the OperatorProfile.
func (OperatorProfile) Fields() []ent.Field {
	return nil
}

// Edges of the OperatorProfile.
func (OperatorProfile) Edges() []ent.Edge {
	return nil
}
