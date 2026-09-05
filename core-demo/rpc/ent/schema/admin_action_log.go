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
	return []schema.Annotation{entsql.Annotation{Table: "sys_admin_action_log"}}
}

func (AdminActionLog) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.OperatorIDMixin{}}
}

func (AdminActionLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id"),
		field.String("request_method").MaxLen(10),
		field.String("request_path").MaxLen(500),
		field.Text("request_query").Optional().Nillable(),
		field.Text("request_body").Optional().Nillable(),
		field.Int16("action_result"),
		field.Int("response_status"),
		field.Text("response_body").Optional().Nillable(),
		field.Int("duration_ms").Default(0),
		field.String("client_ip").
			SchemaType(map[string]string{dialect.Postgres: "inet"}),
		field.String("user_agent").MaxLen(1000).Optional().Nillable(),
		field.Time("created_at").Default(time.Now),
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
		index.Fields("operator_id", "created_at").
			StorageKey("idx_sys_admin_action_log_operator_time"),
	}
}
