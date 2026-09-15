package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// GameCurrency holds the schema definition for the GameCurrency entity.
type GameCurrency struct {
	ent.Schema
}

// Table of the GameCurrency.
func (GameCurrency) Table() string {
	return "game_currency"
}

// Fields of the GameCurrency.
func (GameCurrency) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.Int64("game_id"),
		field.Int64("currency_id"),
		field.Int64("source_status"),
		field.Int64("status").Range(1, 2),
		field.Time("deleted_at").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the GameCurrency.
func (GameCurrency) Edges() []ent.Edge {
	return nil
}

func (GameCurrency) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                 // 启用数据库字段注释
		schema.Comment("系统游戏货币表"),                 // 设置数据库表注释
		entsql.Annotation{Table: "game_currency"}, // 设置数据库表名
	}
}
