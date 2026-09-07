package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type I18n struct{ ent.Schema }

func (I18n) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "sys_i18n"}}
}

func (I18n) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (I18n) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("i18n_group").MaxLen(64),
		field.String("trans_key").MaxLen(192),
		field.String("lang").MaxLen(16),
		field.String("value").MaxLen(1024),
	}
}

func (I18n) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("trans_key", "lang").Unique().StorageKey("uk_sys_i18n_trans_key_lang"),
		index.Fields("i18n_group", "lang").StorageKey("idx_sys_i18n_group_lang"),
	}
}
