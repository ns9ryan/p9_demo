package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type CasbinRule struct{ ent.Schema }

func (CasbinRule) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "casbin_rule"}}
}

func (CasbinRule) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("ptype").MaxLen(16).Default("p"),
		field.String("v0").MaxLen(255).Default(""),
		field.String("v1").MaxLen(255).Default(""),
		field.String("v2").MaxLen(255).Default(""),
		field.String("v3").MaxLen(255).Default(""),
		field.String("v4").MaxLen(255).Default(""),
		field.String("v5").MaxLen(255).Default(""),
	}
}

func (CasbinRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ptype", "v1", "v0", "v2", "v3", "v4", "v5").Unique().
			StorageKey("uk_casbin_rule_policy"),
	}
}
