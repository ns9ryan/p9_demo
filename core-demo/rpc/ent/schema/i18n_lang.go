package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type I18nLang struct{ ent.Schema }

func (I18nLang) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("I18n Language Table | 语言表"),
		entsql.Annotation{Table: "sys_i18n_lang"},
	}
}

func (I18nLang) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (I18nLang) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Primary key | 主键"),
		field.String("lang").MaxLen(16).Comment("Language code | 语言码"),
		field.String("name").MaxLen(64).Default("").Comment("Display name | 显示名"),
		field.String("i18n_key").MaxLen(255).Default("").Comment("I18n key | 多语言 key"),
		field.Int16("disabled").Default(0).Comment("Disabled 0 no 1 yes | 停用 0 否 1 是"),
		field.Int("sort_no").Default(0).Comment("Sort order | 排序"),
	}
}

func (I18nLang) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("lang").Unique().StorageKey("uk_sys_i18n_lang"),
	}
}
