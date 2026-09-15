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
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("I18n Table | 多语言词条表"),
		entsql.Annotation{Table: "sys_i18n"},
	}
}

func (I18n) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (I18n) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Primary key | 主键"),
		field.String("i18n_code").MaxLen(32).Default("platform").Comment("Site code | 站点编码"),
		field.String("i18n_group").MaxLen(64).Comment("Group menu/api/front | 分组"),
		field.String("trans_key").MaxLen(255).Comment("Translation key | 词条 key"),
		field.String("lang").MaxLen(16).Comment("Language code | 语言码"),
		field.String("value").MaxLen(1024).Comment("Translated text | 译文"),
	}
}

func (I18n) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("i18n_code", "trans_key", "lang").Unique().StorageKey("uk_sys_i18n_code_key_lang"),
		index.Fields("i18n_code", "i18n_group", "lang").StorageKey("idx_sys_i18n_code_group_lang"),
	}
}
