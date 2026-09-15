package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// GameProvider holds the schema definition for the GameProvider entity.
type GameProvider struct {
	ent.Schema
}

// Table of the GameProvider.
func (GameProvider) Table() string {
	return "game_provider"
}

// Fields of the GameProvider.
func (GameProvider) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.Int64("source_id"),
		field.String("provider_code"),
		field.String("source_provider_code"),
		field.String("source_logo_url").Optional(),
		field.String("logo_url").Optional(),
		field.Int64("source_sort_no").Optional(),
		field.Int64("sort_no"),
		field.Int64("source_status"),
		field.Int64("status").Range(1, 2),
		field.Time("deleted_at").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the GameProvider.
func (GameProvider) Edges() []ent.Edge {
	return nil
}

func (GameProvider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                 // 启用数据库字段注释
		schema.Comment("系统游戏供应商表"),                // 设置数据库表注释
		entsql.Annotation{Table: "game_provider"}, // 设置数据库表名
	}
}
