package schema

import (
	"oa.98ent.com/p9/node-dispatch/rpc/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// OperatorNode 定义 operator 节点关系表结构
type OperatorNode struct {
	ent.Schema
}

// Fields 定义 operator 节点关系表字段
func (OperatorNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("operator_code").
			NotEmpty().
			MaxLen(64).
			Unique().
			Immutable().
			Comment("operator 全局唯一业务编码"),

		field.Int64("node_id").
			Comment("所属节点本地主键"),
	}
}

// Edges 定义 operator 节点关系表关联关系
func (OperatorNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node", Node.Type).
			Ref("operator_nodes").
			Field("node_id").
			Unique().
			Required(),
	}
}

// Indexes 定义 operator 节点关系表索引
func (OperatorNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id"),
	}
}

// Mixin 定义 operator 节点关系表公共字段
func (OperatorNode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义 operator 节点关系表数据库注解
func (OperatorNode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                 // 启用数据库字段注释
		schema.Comment("operator 节点关系表"),          // 设置数据库表注释
		entsql.Annotation{Table: "operator_node"}, // 设置数据库表名
	}
}
