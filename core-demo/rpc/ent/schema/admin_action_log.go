package schema

import (
	"time"

	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type AdminActionLog struct{ ent.Schema }

func (AdminActionLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("Admin Action Log Table | 操作日志表"),
		entsql.Annotation{Table: "sys_admin_action_log"},
	}
}

func (AdminActionLog) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.IDMixin{}, entmixin.OperatorCodeMixin{}}
}

func (AdminActionLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id").Comment("User ID | 用户 ID"),
		field.String("request_method").MaxLen(10).Comment("HTTP method | 请求方法"),
		field.String("request_path").MaxLen(500).Comment("Request path | 请求路径"),
		field.Text("request_query").Optional().Nillable().Comment("Query string | 查询串"),
		field.Text("request_body").Optional().Nillable().Comment("Request body | 请求体"),
		field.Int16("action_result").Comment("Action result 1 success 2 fail | 操作结果 1 成功 2 失败"),
		field.Int("response_status").Comment("HTTP status | 响应状态码"),
		field.Text("response_body").Optional().Nillable().Comment("Response body | 响应体"),
		field.Int("duration_ms").Default(0).Comment("Duration ms | 耗时毫秒"),
		field.String("client_ip").
			SchemaType(map[string]string{dialect.Postgres: "inet"}).
			Comment("Client IP | 客户端 IP"),
		field.String("user_agent").MaxLen(1000).Optional().Nillable().Comment("User agent | 客户端标识"),
		field.Time("created_at").Default(time.Now).Comment("Created at | 创建时间"),
	}
}

func (AdminActionLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("action_logs").
			Field("user_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Restrict)),
	}
}

func (AdminActionLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at").
			StorageKey("idx_sys_admin_action_log_time"),
		index.Fields("user_id", "created_at").
			StorageKey("idx_sys_admin_action_log_user_time"),
		index.Fields("request_method", "request_path", "created_at").
			StorageKey("idx_sys_admin_action_log_api_time"),
		index.Fields("operator_code", "created_at").
			StorageKey("idx_sys_admin_action_log_operator_time"),
	}
}
