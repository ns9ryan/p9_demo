package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// OperatorGameChannel holds the schema definition for the OperatorGameChannel entity.
type OperatorGameChannel struct {
	ent.Schema
}

// Table of the OperatorGameChannel.
func (OperatorGameChannel) Table() string {
	return "operator_game_channel"
}

// Fields of the OperatorGameChannel.
func (OperatorGameChannel) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("op_code"),
		field.String("channel_code"),
		field.Int16("status").Range(1, 2),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the OperatorGameChannel.
func (OperatorGameChannel) Edges() []ent.Edge {
	return nil
}

func (OperatorGameChannel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                         // 启用数据库字段注释
		schema.Comment("分站游戏渠道关联表"),                       // 设置数据库表注释
		entsql.Annotation{Table: "operator_game_channel"}, // 设置数据库表名
	}
}
