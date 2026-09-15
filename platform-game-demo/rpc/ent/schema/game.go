package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Game holds the schema definition for the Game entity.
type Game struct {
	ent.Schema
}

// Table of the Game.
func (Game) Table() string {
	return "game"
}

// Fields of the Game.
func (Game) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("game_code"),
		field.String("source_game_code"),
		field.Int64("source_id").Optional(),
		field.Int64("category_id"),
		field.Int64("provider_id"),
		field.String("source_provider_key").Optional(),
		field.String("provider_key").Optional(),
		field.Int64("channel_id").Optional(),
		field.String("name").Optional(),
		field.String("source_image_url").Optional(),
		field.String("image_url").Optional(),
		field.Int64("source_sort_no").Optional(),
		field.Int64("sort_no"),
		field.Bool("supports_embed"),
		field.Bool("supports_redirect"),
		field.Int64("source_status"),
		field.Int64("status").Range(1, 2),
		field.Time("deleted_at").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the Game.
func (Game) Edges() []ent.Edge {
	return nil
}

func (Game) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),        // 启用数据库字段注释
		schema.Comment("系统游戏表"),          // 设置数据库表注释
		entsql.Annotation{Table: "game"}, // 设置数据库表名
	}
}
