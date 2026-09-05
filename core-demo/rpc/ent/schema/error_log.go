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
	return []schema.Annotation{entsql.Annotation{Table: "sys_error_log"}}
}

func (ErrorLog) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.OperatorIDMixin{}}
}

func (ErrorLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id").Optional().Nillable(),
		field.String("request_method").MaxLen(10),
		field.String("request_path").MaxLen(500),
		field.Text("request_query").Optional().Nillable(),
		field.Text("request_body").Optional().Nillable(),
		field.String("service_name").MaxLen(100),
		field.Int("response_status"),
		field.Text("response_body").Optional().Nillable(),
		field.Text("subject").Optional().Nillable(),
		field.Text("detail").Optional().Nillable(),
		field.Int("duration_ms").Default(0),
		field.String("client_ip").
			SchemaType(map[string]string{dialect.Postgres: "inet"}),
		field.String("user_agent").MaxLen(1000).Optional().Nillable(),
		field.Time("created_at").Default(time.Now),
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
		index.Fields("operator_id", "created_at").
			StorageKey("idx_sys_error_log_operator_time"),
		index.Fields("service_name", "created_at").
			StorageKey("idx_sys_error_log_service_time"),
	}
}
