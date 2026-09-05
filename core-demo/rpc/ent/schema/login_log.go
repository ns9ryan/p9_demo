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

type LoginLog struct{ ent.Schema }

func (LoginLog) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "sys_login_log"}}
}

func (LoginLog) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.OperatorIDMixin{}}
}

func (LoginLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("user_id").Optional().Nillable(),
		field.String("username").MaxLen(64),
		field.Int16("login_result"),
		field.String("failure_reason").MaxLen(255).Optional().Nillable(),
		field.String("login_ip").
			SchemaType(map[string]string{dialect.Postgres: "inet"}),
		field.Int64("device_id").Optional().Nillable(),
		field.String("user_agent").MaxLen(1000).Optional().Nillable(),
		field.Time("login_at").Default(time.Now),
	}
}

func (LoginLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("login_logs").
			Field("user_id").
			Unique().
			Annotations(entsql.OnDelete(entsql.Restrict)),
	}
}

func (LoginLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "login_at").
			StorageKey("idx_sys_login_log_user_time").
			Annotations(entsql.IndexWhere("user_id IS NOT NULL")),
		index.Fields("login_result", "login_at").
			StorageKey("idx_sys_login_log_result_time"),
		index.Fields("login_at").
			StorageKey("idx_sys_login_log_time"),
		index.Fields("operator_id", "login_at").
			StorageKey("idx_sys_login_log_operator_time"),
	}
}
