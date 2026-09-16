package schema

import (
	"encoding/json"

	"oa.98ent.com/p9/node-dispatch/rpc/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DispatchTask 定义调度任务表结构
type DispatchTask struct {
	ent.Schema
}

// Fields 定义调度任务表字段
func (DispatchTask) Fields() []ent.Field {
	return []ent.Field{
		field.String("task_no").
			NotEmpty().
			MaxLen(64).
			Unique().
			Immutable().
			Comment("全局唯一任务编号"),

		field.String("task_type").
			NotEmpty().
			MaxLen(64).
			Immutable().
			Comment("任务类型"),

		field.JSON("params", json.RawMessage{}).
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Comment("任务参数"),

		field.Int64("status").
			Default(1).
			Range(1, 4).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("任务状态: 1待执行, 2执行中, 3成功, 4失败"),
	}
}

// Edges 定义调度任务表关联关系
func (DispatchTask) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("runs", DispatchTaskRun.Type),
	}
}

// Mixin 定义调度任务表公共字段
func (DispatchTask) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义调度任务表数据库注解
func (DispatchTask) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                 // 启用数据库字段注释
		schema.Comment("调度任务表"),                   // 设置数据库表注释
		entsql.Annotation{Table: "dispatch_task"}, // 设置数据库表名
	}
}
