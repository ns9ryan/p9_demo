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
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("User Table | 用户表"),
		entsql.Annotation{Table: "sys_user"},
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{entmixin.TimeMixin{}, entmixin.SoftDeleteMixin{}, entmixin.OperatorCodeMixin{}}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Primary key | 主键"),
		field.String("user_code").MaxLen(64).Comment("User code | 用户编码"),
		field.String("username").MaxLen(64).Comment("Login name | 登录名"),
		field.String("password_hash").MaxLen(255).Sensitive().Comment("Password hash | 密码哈希"),
		field.String("salt").MaxLen(64).Sensitive().Comment("Token salt | 令牌盐"),
		field.String("display_name").MaxLen(100).Comment("Display name | 显示名"),
		field.String("mobile").MaxLen(32).Optional().Nillable().Comment("Mobile | 手机号"),
		field.String("email").MaxLen(255).Optional().Nillable().Comment("Email | 邮箱"),
		field.Int16("status").Default(1).Comment("Status 1 enabled 2 disabled | 状态 1 启用 2 停用"),
		field.Bool("is_super_admin").Default(false).Comment("Super admin | 是否超管"),
		field.Time("last_login_at").Optional().Nillable().Comment("Last login time | 最后登录时间"),
		field.String("last_login_ip").Optional().Nillable().
			SchemaType(map[string]string{dialect.Postgres: "inet"}).
			Comment("Last login IP | 最后登录 IP"),
		field.Int16("ip_whitelist_enabled").Default(0).
			Comment("IP whitelist switch 0 off 1 on | IP 白名单开关 0 关 1 开"),
		field.JSON("ip_whitelist", []string{}).
			Default([]string{}).
			SchemaType(map[string]string{dialect.Postgres: "jsonb"}).
			Annotations(entsql.DefaultExpr("'[]'")).
			Comment("Allowed login IPs or CIDRs | 允许登录的 IP / CIDR"),
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
		index.Fields("operator_code", "username").Unique().
			StorageKey("uk_sys_user_operator_username").
			Annotations(entsql.IndexWhere("operator_code IS NOT NULL")),
		index.Fields("username").Unique().
			StorageKey("uk_sys_user_username").
			Annotations(entsql.IndexWhere("operator_code IS NULL")),
		index.Fields("operator_code").Unique().
			StorageKey("uk_sys_user_operator_super_admin").
			Annotations(entsql.IndexWhere("is_super_admin IS TRUE AND deleted_at IS NULL")),
		index.Fields("operator_code").
			StorageKey("idx_sys_user_operator"),
	}
}
