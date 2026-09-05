package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Promotion struct{ ent.Schema }

func (Promotion) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "partner_promotion"}}
}

func (Promotion) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}, entmixin.OperatorCodeMixin{}}
}

func (Promotion) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("title").MaxLen(255),
	}
}
