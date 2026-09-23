package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"oa.98ent.com/p9/operator-base/rpc/ent/schema/mixins"
)

// Operator 定义厅 operator 基础表结构
type Operator struct {
	ent.Schema
}

// Fields 定义厅 operator 基础表字段
func (Operator) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty().
			MaxLen(64).
			Unique().
			Immutable().
			Comment("operator 全局唯一业务编码"),

		field.String("name").
			NotEmpty().
			MaxLen(100).
			Comment("operator 名称"),

		field.String("timezone_code").
			NotEmpty().
			MaxLen(64).
			Immutable().
			Comment("IANA 时区编码"),

		field.String("settlement_currency_code").
			NotEmpty().
			MaxLen(16).
			Immutable().
			Comment("结算货币编码"),

		field.Int64("status").
			Default(1).
			Range(1, 3).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("operator 状态: 1正常, 2暂停, 3关闭"),
	}
}

// Edges 定义厅 operator 基础表关联关系
func (Operator) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("domains", OperatorDomain.Type),
		edge.To("languages", OperatorLanguage.Type),
		edge.To("regions", OperatorRegion.Type),
		edge.To("agent_lines", OperatorAgentLine.Type),
	}
}

// Mixin 定义厅 operator 基础表公共字段
func (Operator) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义厅 operator 基础表数据库注解
func (Operator) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),            // 启用数据库字段注释
		schema.Comment("厅 operator 基础表"),     // 设置数据库表注释
		entsql.Annotation{Table: "operator"}, // 设置数据库表名
	}
}
