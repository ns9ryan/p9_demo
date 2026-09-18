package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// OperatorGameProvider holds the schema definition for the OperatorGameProvider entity.
type OperatorGameProvider struct {
	ent.Schema
}

// Table of the OperatorGameProvider.
func (OperatorGameProvider) Table() string {
	return "operator_game_provider"
}

// Fields of the OperatorGameProvider.
func (OperatorGameProvider) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("op_code"),
		field.String("provider_code"),
		field.Int16("status").Range(1, 2),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the OperatorGameProvider.
func (OperatorGameProvider) Edges() []ent.Edge {
	return nil
}

func (OperatorGameProvider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                          // 启用数据库字段注释
		schema.Comment("分站游戏提供商关联表"),                       // 设置数据库表注释
		entsql.Annotation{Table: "operator_game_provider"}, // 设置数据库表名
	}
}
