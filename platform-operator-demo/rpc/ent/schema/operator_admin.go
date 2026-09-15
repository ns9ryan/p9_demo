package schema

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// OperatorAdmin 定义 operator 管理员表结构
type OperatorAdmin struct {
	ent.Schema
}

// Fields 定义 operator 管理员表字段
func (OperatorAdmin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("operator_id").
			Comment("所属 operator 本地主键"),

		field.String("username").
			NotEmpty().
			MaxLen(64).
			Comment("账号"),

		field.String("password").
			NotEmpty().
			MaxLen(64).
			Sensitive().
			Comment("密码"),

		field.String("display_name").
			NotEmpty().
			MaxLen(100).
			Comment("显示名称"),
	}
}

// Edges 定义 operator 管理员表关联关系
func (OperatorAdmin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", Operator.Type).
			Ref("admins").
			Field("operator_id").
			Unique().
			Required(),
	}
}

// Indexes 定义 operator 管理员表索引
func (OperatorAdmin) Indexes() []ent.Index {
	return []ent.Index{
		// 同一个 operator 下账号唯一
		index.Fields("operator_id", "username").
			Unique().
			StorageKey("uk_operator_admin_operator_username"),
	}
}

// Mixin 定义 operator 管理员表公共字段
func (OperatorAdmin) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义 operator 管理员表数据库注解
func (OperatorAdmin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                  // 启用数据库字段注释
		schema.Comment("总网 operator 管理员表"),         // 设置数据库表注释
		entsql.Annotation{Table: "operator_admin"}, // 设置数据库表名
	}
}
