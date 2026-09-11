package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"oa.98ent.com/p9/platform-operator/rpc/ent/schema/mixins"
)

// OperatorDomain 定义 operator 域名表结构
type OperatorDomain struct {
	ent.Schema
}

// Fields 定义 operator 域名表字段
func (OperatorDomain) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("operator_id").
			Comment("所属 operator 本地主键"),

		field.String("domain_name").
			NotEmpty().
			MaxLen(253).
			Comment("域名"),

		field.Int64("domain_type").
			Range(1, 3).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("域名类型: 1分站后台, 2代理后台, 3会员H5"),

		field.String("remark").
			MaxLen(1000).
			Optional().
			Nillable().
			Comment("总网内部备注"),
	}
}

// Edges 定义 operator 域名表关联关系
func (OperatorDomain) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", Operator.Type).
			Ref("domains").
			Field("operator_id").
			Unique().
			Required(),
	}
}

// Indexes 定义 operator 域名表索引
func (OperatorDomain) Indexes() []ent.Index {
	return []ent.Index{
		// 同一个 operator 不能重复保存同一个域名
		index.Fields("operator_id", "domain_name").
			Unique().
			StorageKey("uk_operator_domain_operator_name"),

		// 同一个启用中的域名不能同时绑定多个 operator
		index.Fields("domain_name").
			Unique().
			StorageKey("uk_operator_domain_active_name").
			Annotations(
				entsql.IndexWhere("status = 1"),
			),

		// 同一个 operator 的同一种域名类型只能有一个启用中的域名
		index.Fields("operator_id", "domain_type").
			Unique().
			StorageKey("uk_operator_domain_active_type").
			Annotations(
				entsql.IndexWhere("status = 1"),
			),
	}
}

// Mixin 定义 operator 域名表公共字段
func (OperatorDomain) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
		mixins.TimeMixin{},
	}
}

// Annotations 定义 operator 域名表数据库注解
func (OperatorDomain) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                   // 启用数据库字段注释
		schema.Comment("总网 operator 域名表"),           // 设置数据库表注释
		entsql.Annotation{Table: "operator_domain"}, // 设置数据库表名
	}
}
