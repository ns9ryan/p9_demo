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
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("API Table | 接口表"),
		entsql.Annotation{Table: "sys_api"},
	}
}

func (API) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (API) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Primary key | 主键"),
		field.String("description").MaxLen(255).Comment("Description i18n key | 描述词条"),
		field.String("api_group").MaxLen(100).Comment("API group | 接口分组"),
		field.String("method").MaxLen(10).Comment("HTTP method | 请求方法"),
		field.String("path").MaxLen(255).Comment("Request path | 请求路径"),
		field.Int16("is_required").Default(0).Comment("Required 0 no 1 yes | 是否必选 0 否 1 是"),
		field.String("service_name").MaxLen(255).Comment("Service name | 服务名"),
	}
}

func (API) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("method", "path").Unique().StorageKey("uk_sys_api_method_path"),
	}
}
