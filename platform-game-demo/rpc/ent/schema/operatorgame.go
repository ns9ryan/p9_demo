package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// OperatorGame holds the schema definition for the OperatorGame entity.
type OperatorGame struct {
	ent.Schema
}

// Table of the OperatorGame.
func (OperatorGame) Table() string {
	return "operator_game"
}

// Fields of the OperatorGame.
func (OperatorGame) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("op_code"),
		field.String("game_code"),
		field.String("name"),
		field.Int16("status").Range(1, 2),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the OperatorGame.
func (OperatorGame) Edges() []ent.Edge {
	return nil
}

func (OperatorGame) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                 // 启用数据库字段注释
		schema.Comment("分站游戏关联表"),                 // 设置数据库表注释
		entsql.Annotation{Table: "operator_game"}, // 设置数据库表名
	}
}
