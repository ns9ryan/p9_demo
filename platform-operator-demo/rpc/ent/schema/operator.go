package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"oa.98ent.com/p9/platform-operator/rpc/ent/schema/mixins"
)

// Operator 定义 operator 基础表结构
type Operator struct {
	ent.Schema
}

// Fields 定义 operator 基础表字段
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
			Comment("IANA 时区编码"),

		field.String("settlement_currency_code").
			NotEmpty().
			MaxLen(16).
			Comment("结算货币编码"),

		field.Int64("creation_status").
			Default(1).
			Range(1, 2).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("创建状态: 1草稿, 2已完成"),

		field.Int64("publish_status").
			Default(1).
			Range(1, 4).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("发布状态: 1未发布, 2发布中, 3已发布, 4发布失败"),

		field.String("publish_task_no").
			MaxLen(64).
			Optional().
			Nillable().
			Unique().
			Comment("最近一次发布任务编号"),

		field.Time("published_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz(3)",
			}).
			Comment("发布时间"),

		field.Int64("status").
			Default(1).
			Range(1, 3).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("operator 状态: 1正常, 2暂停, 3关闭"),

		field.String("remark").
			MaxLen(1000).
			Optional().
			Nillable().
			Comment("总网内部备注"),
	}
}

// Edges 定义 operator 基础表关联关系
func (Operator) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("profile", OperatorProfile.Type).
			Unique(),
		edge.To("domains", OperatorDomain.Type),
		edge.To("language_allocations", OperatorLanguageAllocation.Type),
		edge.To("region_allocations", OperatorRegionAllocation.Type),
		edge.To("agent_line_allocations", OperatorAgentLineAllocation.Type),
	}
}

// Mixin 定义 operator 基础表公共字段
func (Operator) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义 operator 基础表数据库注解
func (Operator) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),            // 启用数据库字段注释
		schema.Comment("总网 operator 基础表"),    // 设置数据库表注释
		entsql.Annotation{Table: "operator"}, // 设置数据库表名
	}
}
