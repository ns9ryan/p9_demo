package schema

import (
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"oa.98ent.com/p9/platform-operator/rpc/ent/schema/mixins"
)

// OperatorRegionAllocation 定义 operator 经营地区分配表结构
type OperatorRegionAllocation struct {
	ent.Schema
}

// Fields 定义 operator 经营地区分配表字段
func (OperatorRegionAllocation) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("operator_id").
			Comment("所属 operator 本地主键"),

		field.String("region_code").
			NotEmpty().
			MinLen(2).
			MaxLen(2).
			Match(regexp.MustCompile(`^[A-Z]{2}$`)).
			Comment("国家或地区唯一业务编码"),

		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz(3)",
			}).
			Comment("创建时间"),
	}
}

// Edges 定义 operator 经营地区分配表关联关系
func (OperatorRegionAllocation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", Operator.Type).
			Ref("region_allocations").
			Field("operator_id").
			Unique().
			Required(),
	}
}

// Indexes 定义 operator 经营地区分配表索引
func (OperatorRegionAllocation) Indexes() []ent.Index {
	return []ent.Index{
		// 同一个 operator 与同一个经营地区只保存一条分配关系
		index.Fields("operator_id", "region_code").
			Unique().
			StorageKey("uk_operator_region_allocation"),
	}
}

// Mixin 定义 operator 经营地区分配表公共字段
func (OperatorRegionAllocation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
	}
}

// Annotations 定义 operator 经营地区分配表数据库注解
func (OperatorRegionAllocation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                              // 启用数据库字段注释
		schema.Comment("总网 operator 经营地区分配表"),                  // 设置数据库表注释
		entsql.Annotation{Table: "operator_region_allocation"}, // 设置数据库表名
	}
}
