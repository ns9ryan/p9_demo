package schema

import "entgo.io/ent"

// Operator holds the schema definition for the Operator entity.
type Operator struct {
	ent.Schema
}

// Fields of the Operator.
func (Operator) Fields() []ent.Field {
	return nil
}

// Edges of the Operator.
func (Operator) Edges() []ent.Edge {
	return nil
}
