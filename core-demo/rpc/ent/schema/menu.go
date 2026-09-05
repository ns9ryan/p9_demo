package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Menu struct{ ent.Schema }

func (Menu) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "sys_menu"}}
}

func (Menu) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}}
}

func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("parent_id").Default(0),
		field.Int16("menu_type"),
		field.String("path").MaxLen(128).Default(""),
		field.String("name").MaxLen(64).Default(""),
		field.String("component").MaxLen(255).Default(""),
		field.String("redirect").MaxLen(255).Default(""),
		field.String("title").MaxLen(64).Default(""),
		field.String("icon").MaxLen(64).Default(""),
		field.String("permission").MaxLen(128).Default(""),
		field.Int16("hide_menu").Default(0),
		field.Int("sort").Default(0),
		field.Int16("disabled").Default(0),
	}
}

func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", Role.Type).Ref("menus"),
	}
}

func (Menu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique().StorageKey("uk_sys_menu_name"),
	}
}
