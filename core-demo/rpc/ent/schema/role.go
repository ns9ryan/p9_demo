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

type Role struct{ ent.Schema }

func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "sys_role"}}
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}, entmixin.SoftDeleteMixin{}, entmixin.OperatorIDMixin{}}
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("role_code").MaxLen(64),
		field.String("role_name").MaxLen(100),
		field.String("description").MaxLen(255).Optional().Nillable(),
		field.Int16("status").Default(1),
		field.Bool("is_system").Default(false),
		field.Int("sort_no").Default(0),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("roles"),
		edge.To("menus", Menu.Type).
			StorageKey(edge.Table("sys_role_menu"), edge.Columns("role_id", "menu_id")),
	}
}

func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("operator_id", "role_code").Unique().
			StorageKey("uk_sys_role_operator_code").
			Annotations(entsql.IndexWhere("operator_id IS NOT NULL")),
		index.Fields("role_code").Unique().
			StorageKey("uk_sys_role_code").
			Annotations(entsql.IndexWhere("operator_id IS NULL")),
		index.Fields("operator_id", "role_name").Unique().
			StorageKey("uk_sys_role_operator_name").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
	}
}
