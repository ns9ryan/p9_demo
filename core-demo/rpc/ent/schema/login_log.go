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
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("Login Log Table | 登录日志表"),
		entsql.Annotation{Table: "sys_login_log"},
	}
}

func (LoginLog) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.IDMixin{}, entmixin.OperatorCodeMixin{}}
}

func (LoginLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id").Optional().Nillable().Comment("User ID | 用户 ID"),
		field.String("username").MaxLen(64).Comment("Login name | 登录名"),
		field.Int16("login_result").Comment("Login result 1 success 2 fail | 登录结果 1 成功 2 失败"),
		field.String("failure_reason").MaxLen(255).Optional().Nillable().Comment("Failure reason | 失败原因"),
		field.String("login_ip").
			SchemaType(map[string]string{dialect.Postgres: "inet"}).
			Comment("Login IP | 登录 IP"),
		field.Int64("device_id").Optional().Nillable().Comment("Device ID | 设备 ID"),
		field.String("user_agent").MaxLen(1000).Optional().Nillable().Comment("User agent | 客户端标识"),
		field.Time("login_at").Default(time.Now).Comment("Login time | 登录时间"),
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
		index.Fields("operator_code", "login_at").
			StorageKey("idx_sys_login_log_operator_time"),
	}
}
