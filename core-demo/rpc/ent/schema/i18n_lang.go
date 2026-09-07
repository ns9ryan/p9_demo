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
	return []schema.Annotation{entsql.Annotation{Table: "sys_i18n_lang"}}
}

func (I18nLang) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (I18nLang) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("lang").MaxLen(16),
		field.String("name").MaxLen(64).Default(""),
		field.Int16("disabled").Default(0),
		field.Int16("is_default").Default(0),
	}
}

func (I18nLang) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("lang").Unique().StorageKey("uk_sys_i18n_lang"),
	}
}
