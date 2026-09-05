package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type API struct{ ent.Schema }

func (API) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "sys_api"}}
}

func (API) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (API) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("description").MaxLen(255),
		field.String("api_group").MaxLen(100),
		field.String("method").MaxLen(10),
		field.String("path").MaxLen(255),
		field.Int16("is_required").Default(0),
		field.String("service_name").MaxLen(255),
	}
}

func (API) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("method", "path").Unique().StorageKey("uk_sys_api_method_path"),
	}
}
