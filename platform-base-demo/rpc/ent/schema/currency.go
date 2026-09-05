package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"oa.98ent.com/p9/platform-base/rpc/ent/schema/mixins"
)

// Currency 定义系统货币表结构
type Currency struct {
	ent.Schema
}

// Fields 定义系统货币表字段
func (Currency) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty().
			MaxLen(16).
			Unique().
			Immutable().
			Comment("货币编码"),

		field.JSON("name_i18n", map[string]string{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Comment("多语言名称"),

		field.Int64("currency_type").
			Range(1, 2).
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("货币类型: 1法定货币, 2虚拟货币"),

		field.String("symbol").
			NotEmpty().
			MaxLen(16).
			Comment("货币符号"),

		field.Int64("amount_factor").
			Min(1).
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "bigint",
			}).
			Comment("金额换算倍率, 例如 USD 为100, VND为1"),
	}
}

// Edges 定义系统货币表关联关系
func (Currency) Edges() []ent.Edge {
	return nil
}

// Mixin 定义系统货币表公共字段
func (Currency) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.SortMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义系统货币表数据库注解
func (Currency) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),            // 启用数据库字段注释
		schema.Comment("系统货币表"),              // 设置数据库表注释
		entsql.Annotation{Table: "currency"}, // 设置数据库表名
	}
}
