package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SysInit struct{ ent.Schema }

func (SysInit) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("Catalog init flag | 目录初始化标记"),
		entsql.Annotation{Table: "sys_init"},
	}
}

func (SysInit) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.IDMixin{}, entmixin.TimeMixin{}}
}

func (SysInit) Fields() []ent.Field {
	return []ent.Field{
		field.String("init_key").MaxLen(32).Comment("Init category key | 初始化类别"),
	}
}

func (SysInit) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("init_key").Unique().StorageKey("uk_sys_init_key"),
	}
}
