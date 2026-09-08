package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"oa.98ent.com/p9/platform-operator/rpc/ent/schema/mixins"
)

// OperatorProfile 定义 operator 档案表结构
type OperatorProfile struct {
	ent.Schema
}

// Fields 定义 operator 档案表字段
func (OperatorProfile) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("operator_id").
			Unique().
			Comment("所属 operator 本地主键"),

		field.String("company_name").
			MaxLen(200).
			Optional().
			Nillable().
			Comment("公司名称"),

		field.String("contact_name").
			MaxLen(100).
			Optional().
			Nillable().
			Comment("主要联系人名称"),

		field.String("contact_email").
			MaxLen(255).
			Optional().
			Nillable().
			Comment("主要联系人邮箱"),

		field.String("remark").
			MaxLen(1000).
			Optional().
			Nillable().
			Comment("总网内部档案备注"),
	}
}

// Edges 定义 operator 档案表关联关系
func (OperatorProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", Operator.Type).
			Ref("profile").
			Field("operator_id").
			Unique().
			Required(),
	}
}

// Mixin 定义 operator 档案表公共字段
func (OperatorProfile) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义 operator 档案表数据库注解
func (OperatorProfile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                    // 启用数据库字段注释
		schema.Comment("总网 operator 档案表"),            // 设置数据库表注释
		entsql.Annotation{Table: "operator_profile"}, // 设置数据库表名
	}
}
