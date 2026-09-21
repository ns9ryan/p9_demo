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

type ErrorLog struct{ ent.Schema }

func (ErrorLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("Error Log Table | 错误日志表"),
		entsql.Annotation{Table: "sys_error_log"},
	}
}

func (ErrorLog) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.IDMixin{}, entmixin.OperatorCodeMixin{}}
}

func (ErrorLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id").Optional().Nillable().Comment("User ID | 用户 ID"),
		field.String("request_method").MaxLen(10).Comment("HTTP method | 请求方法"),
		field.String("request_path").MaxLen(500).Comment("Request path | 请求路径"),
		field.Text("request_query").Optional().Nillable().Comment("Query string | 查询串"),
		field.Text("request_body").Optional().Nillable().Comment("Request body | 请求体"),
		field.String("service_name").MaxLen(100).Comment("Service name | 服务名"),
		field.Int("response_status").Comment("HTTP status | 响应状态码"),
		field.Text("response_body").Optional().Nillable().Comment("Response body | 响应体"),
		field.Text("subject").Optional().Nillable().Comment("Error subject | 错误摘要"),
		field.Text("detail").Optional().Nillable().Comment("Error detail | 错误详情"),
		field.Int("duration_ms").Default(0).Comment("Duration ms | 耗时毫秒"),
		field.String("client_ip").
			SchemaType(map[string]string{dialect.Postgres: "inet"}).
			Comment("Client IP | 客户端 IP"),
		field.String("user_agent").MaxLen(1000).Optional().Nillable().Comment("User agent | 客户端标识"),
		field.Time("created_at").Default(time.Now).Comment("Created at | 创建时间"),
	}
}

func (ErrorLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("error_logs").
			Field("user_id").
			Unique().
			Annotations(entsql.OnDelete(entsql.Restrict)),
	}
}

func (ErrorLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at").
			StorageKey("idx_sys_error_log_time"),
		index.Fields("user_id", "created_at").
			StorageKey("idx_sys_error_log_user_time").
			Annotations(entsql.IndexWhere("user_id IS NOT NULL")),
		index.Fields("request_path", "created_at").
			StorageKey("idx_sys_error_log_path_time"),
		index.Fields("operator_code", "created_at").
			StorageKey("idx_sys_error_log_operator_time"),
		index.Fields("service_name", "created_at").
			StorageKey("idx_sys_error_log_service_time"),
	}
}
