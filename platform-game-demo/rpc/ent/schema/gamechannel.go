package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// GameChannel holds the schema definition for the GameChannel entity.
type GameChannel struct {
	ent.Schema
}

// Table of the GameChannel.
func (GameChannel) Table() string {
	return "game_channel"
}

// Fields of the GameChannel.
func (GameChannel) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.Int64("source_id"),
		field.String("channel_code"),
		field.String("source_channel_code"),
		field.Int64("source_sort_no"),
		field.Int64("sort_no"),
		field.Int64("source_load_type"),
		field.Int64("load_type"),
		field.Int64("source_status"),
		field.Int64("status").Range(1, 2),
		field.Time("deleted_at").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the GameChannel.
func (GameChannel) Edges() []ent.Edge {
	return nil
}

func (GameChannel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                // 启用数据库字段注释
		schema.Comment("系统游戏渠道表"),                // 设置数据库表注释
		entsql.Annotation{Table: "game_channel"}, // 设置数据库表名
	}
}
