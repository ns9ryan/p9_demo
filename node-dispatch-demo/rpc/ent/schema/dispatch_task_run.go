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

// DispatchTaskRun 定义调度任务执行表结构
type DispatchTaskRun struct {
	ent.Schema
}

// Fields 定义调度任务执行表字段
func (DispatchTaskRun) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("task_id").
			Immutable().
			Comment("所属调度任务本地主键"),

		field.Int64("node_id").
			Immutable().
			Comment("执行节点本地主键"),

		field.Int64("status").
			Default(1).
			Range(1, 4).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("执行状态: 1待执行, 2执行中, 3成功, 4失败"),

		field.JSON("result", json.RawMessage{}).
			Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Comment("执行结果"),

		field.String("error_message").
			MaxLen(2000).
			Optional().
			Nillable().
			Comment("执行失败原因"),

		field.Time("started_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz(3)",
			}).
			Comment("开始执行时间"),

		field.Time("finished_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz(3)",
			}).
			Comment("执行结束时间"),
	}
}

// Edges 定义调度任务执行表关联关系
func (DispatchTaskRun) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", DispatchTask.Type).
			Ref("runs").
			Field("task_id").
			Unique().
			Required().
			Immutable(),

		edge.From("node", Node.Type).
			Ref("task_runs").
			Field("node_id").
			Unique().
			Required().
			Immutable(),
	}
}

// Mixin 定义调度任务执行表公共字段
func (DispatchTaskRun) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义调度任务执行表数据库注解
func (DispatchTaskRun) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                     // 启用数据库字段注释
		schema.Comment("调度任务执行表"),                     // 设置数据库表注释
		entsql.Annotation{Table: "dispatch_task_run"}, // 设置数据库表名
	}
}
