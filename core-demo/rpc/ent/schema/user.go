package schema

import (
	"oa.98ent.com/p9/core/common/entmixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type User struct{ ent.Schema }

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "sys_user"}}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}, entmixin.SoftDeleteMixin{}, entmixin.OperatorIDMixin{}}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("user_code").MaxLen(64),
		field.String("username").MaxLen(64),
		field.String("password_hash").MaxLen(255).Sensitive(),
		field.String("salt").MaxLen(64).Sensitive(),
		field.String("display_name").MaxLen(100),
		field.String("mobile").MaxLen(32).Optional().Nillable(),
		field.String("email").MaxLen(255).Optional().Nillable(),
		field.Int16("status").Default(1),
		field.Bool("is_super_admin").Default(false),
		field.Time("last_login_at").Optional().Nillable(),
		field.String("last_login_ip").Optional().Nillable().
			SchemaType(map[string]string{dialect.Postgres: "inet"}),
		field.Int16("ip_whitelist_enabled").Default(0),
		field.JSON("ip_whitelist", []string{}).
			Default([]string{}).
			SchemaType(map[string]string{dialect.Postgres: "jsonb"}).
			Annotations(entsql.DefaultExpr("'[]'")),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", Role.Type).
			StorageKey(edge.Table("sys_user_role"), edge.Columns("user_id", "role_id")),
		edge.To("login_logs", LoginLog.Type),
		edge.To("action_logs", AdminActionLog.Type),
		edge.To("error_logs", ErrorLog.Type),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_code").Unique().
			StorageKey("uk_sys_user_user_code"),
		index.Fields("operator_id", "username").Unique().
			StorageKey("uk_sys_user_operator_username").
			Annotations(entsql.IndexWhere("operator_id IS NOT NULL")),
		index.Fields("username").Unique().
			StorageKey("uk_sys_user_username").
			Annotations(entsql.IndexWhere("operator_id IS NULL")),
		index.Fields("operator_id").Unique().
			StorageKey("uk_sys_user_operator_super_admin").
			Annotations(entsql.IndexWhere("is_super_admin IS TRUE AND deleted_at IS NULL")),
		index.Fields("operator_id").
			StorageKey("idx_sys_user_operator"),
	}
}
