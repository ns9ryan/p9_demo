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
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("Menu Table | 菜单表"),
		entsql.Annotation{Table: "sys_menu"},
	}
}

func (Menu) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.IDMixin{}, entmixin.TimeMixin{}}
}

func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("parent_id").Default(0).Comment("Parent ID | 父级 ID"),
		field.Int16("menu_type").Comment("Menu type 0 directory 1 menu 2 button | 类型 0 目录 1 菜单 2 按钮"),
		field.String("path").MaxLen(128).Default("").Comment("Route path | 路由路径"),
		field.String("name").MaxLen(64).Default("").Comment("Route name | 路由名"),
		field.String("component").MaxLen(255).Default("").Comment("Component | 组件"),
		field.String("redirect").MaxLen(255).Default("").Comment("Redirect | 重定向"),
		field.String("title").MaxLen(64).Default("").Comment("Title i18n key | 标题词条"),
		field.String("icon").MaxLen(64).Default("").Comment("Icon | 图标"),
		field.String("permission").MaxLen(128).Default("").Comment("Permission code | 权限码"),
		field.Int16("hide_menu").Default(0).Comment("Hide menu 0 no 1 yes | 隐藏菜单 0 否 1 是"),
		field.Int("sort").Default(0).Comment("Sort order | 排序"),
		field.Int16("disabled").Default(0).Comment("Disabled 0 no 1 yes | 停用 0 否 1 是"),
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
