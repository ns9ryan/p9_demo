package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Operator struct{ ent.Schema }

func (Operator) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "operator"}}
}

func (Operator) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (Operator) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("operator_code").MaxLen(64),
		field.String("timezone_code").MaxLen(64),
		field.String("settlement_currency_code").MaxLen(16),
		field.Int16("status").Default(1),
		field.Int("required_config_version").Default(1),
		field.Int("completed_config_version").Default(0),
		field.Time("config_completed_at").Optional().Nillable(),
	}
}

func (Operator) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("operator_code").Unique().StorageKey("uk_operator_operator_code"),
	}
}
