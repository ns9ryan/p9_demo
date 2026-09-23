package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"oa.98ent.com/p9/operator-base/rpc/ent/schema/mixins"
)

// OperatorLanguage 定义厅 operator 当前有效语言关系表结构
type OperatorLanguage struct {
	ent.Schema
}

// Fields 定义厅 operator 当前有效语言关系表字段
func (OperatorLanguage) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("operator_id").
			Immutable().
			Comment("所属 operator 本地主键"),

		field.String("language_code").
			NotEmpty().
			MaxLen(16).
			Immutable().
			Comment("系统语言唯一业务编码"),

		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz(3)",
			}).
			Comment("创建时间"),
	}
}

// Edges 定义厅 operator 当前有效语言关系表关联关系
func (OperatorLanguage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", Operator.Type).
			Ref("languages").
			Field("operator_id").
			Unique().
			Required().
			Immutable(),
	}
}

// Indexes 定义厅 operator 当前有效语言关系表索引
func (OperatorLanguage) Indexes() []ent.Index {
	return []ent.Index{
		// 同一个 operator 与同一种语言只保存一条当前有效关系
		index.Fields("operator_id", "language_code").
			Unique().
			StorageKey("uk_operator_language"),
	}
}

// Mixin 定义厅 operator 当前有效语言关系表公共字段
func (OperatorLanguage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
	}
}

// Annotations 定义厅 operator 当前有效语言关系表数据库注解
func (OperatorLanguage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                     // 启用数据库字段注释
		schema.Comment("厅 operator 当前有效语言关系表"),        // 设置数据库表注释
		entsql.Annotation{Table: "operator_language"}, // 设置数据库表名
	}
}
