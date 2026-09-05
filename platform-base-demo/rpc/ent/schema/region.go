package schema

import (
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"oa.98ent.com/p9/platform-base/rpc/ent/schema/mixins"
)

// Region 定义系统国家地区表结构
type Region struct {
	ent.Schema
}

// Fields 定义系统国家地区表字段
func (Region) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty().
			MinLen(2).
			MaxLen(2).
			Match(regexp.MustCompile(`^[A-Z]{2}$`)).
			Unique().
			Immutable().
			Comment("国家或地区编码"),

		field.String("calling_code").
			NotEmpty().
			MaxLen(3).
			Match(regexp.MustCompile(`^[0-9]{1,3}$`)).
			Comment("国际电话区号, 不包含加号"),

		field.JSON("name_i18n", map[string]string{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Comment("多语言名称"),
	}
}

// Edges 定义系统国家地区表关联关系
func (Region) Edges() []ent.Edge {
	return nil
}

// Mixin 定义系统国家地区表公共字段
func (Region) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.SortMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义系统国家地区表数据库注解
func (Region) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),          // 启用数据库字段注释
		schema.Comment("系统国家地区表"),          // 设置数据库表注释
		entsql.Annotation{Table: "region"}, // 设置数据库表名
	}
}
