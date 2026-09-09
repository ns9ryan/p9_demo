package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"oa.98ent.com/p9/platform-base/rpc/ent/schema/mixins"
)

// Timezone 定义系统时区表结构
type Timezone struct {
	ent.Schema
}

// Fields 定义系统时区表字段
func (Timezone) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty().
			MaxLen(64).
			Unique().
			Immutable().
			Comment("IANA 时区编码"),

		field.String("name_key").
			NotEmpty().
			MaxLen(128).
			Immutable().
			Comment("名称翻译 Key"),
	}
}

// Edges 定义系统时区表关联关系
func (Timezone) Edges() []ent.Edge {
	return nil
}

// Mixin 定义系统时区表公共字段
func (Timezone) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.SortMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义系统时区表数据库注解
func (Timezone) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),            // 启用数据库字段注释
		schema.Comment("系统时区表"),              // 设置数据库表注释
		entsql.Annotation{Table: "timezone"}, // 设置数据库表名
	}
}
