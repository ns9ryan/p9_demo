package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// OperatorGameCategory holds the schema definition for the OperatorGameCategory entity.
type OperatorGameCategory struct {
	ent.Schema
}

// Table of the OperatorGameCategory.
func (OperatorGameCategory) Table() string {
	return "operator_game_category"
}

// Fields of the OperatorGameCategory.
func (OperatorGameCategory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("op_code"),
		field.String("category_code"),
		field.Int16("status").Range(1, 2),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the OperatorGameCategory.
func (OperatorGameCategory) Edges() []ent.Edge {
	return nil
}

func (OperatorGameCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                          // 启用数据库字段注释
		schema.Comment("分站游戏分类关联表"),                        // 设置数据库表注释
		entsql.Annotation{Table: "operator_game_category"}, // 设置数据库表名
	}
}
