package schema

import (
	"oa.98ent.com/p9/node-dispatch/rpc/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Node 定义节点表结构
type Node struct {
	ent.Schema
}

// Fields 定义节点表字段
func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty().
			MaxLen(64).
			Unique().
			Immutable().
			Comment("节点全局唯一业务编码"),

		field.String("name").
			NotEmpty().
			MaxLen(100).
			Comment("节点名称"),

		field.String("auth_secret_hash").
			NotEmpty().
			MaxLen(255).
			Sensitive().
			Comment("节点认证密钥哈希"),

		field.Time("last_seen_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz(3)",
			}).
			Comment("最近一次活动时间"),

		field.String("remark").
			MaxLen(1000).
			Optional().
			Nillable().
			Comment("运维备注"),
	}
}

// Edges 定义节点表关联关系
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("operator_nodes", OperatorNode.Type),
		edge.To("task_runs", DispatchTaskRun.Type),
	}
}

// Mixin 定义节点表公共字段
func (Node) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义节点表数据库注解
func (Node) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),        // 启用数据库字段注释
		schema.Comment("节点表"),            // 设置数据库表注释
		entsql.Annotation{Table: "node"}, // 设置数据库表名
	}
}
