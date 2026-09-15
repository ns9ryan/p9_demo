package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// GameCategory holds the schema definition for the GameCategory entity.
type GameCategory struct {
	ent.Schema
}

// Table of the GameCategory.
func (GameCategory) Table() string {
	return "game_category"
}

// Fields of the GameCategory.
func (GameCategory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.Int64("source_id"),
		field.String("category_code"),
		field.String("source_category_code"),
		field.Int64("source_sort_no").Optional(),
		field.Int64("sort_no"),
		field.Int64("source_status"),
		field.Int64("status").Range(1, 2),
		field.Time("deleted_at").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the GameCategory.
func (GameCategory) Edges() []ent.Edge {
	return nil
}

func (GameCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                 // 启用数据库字段注释
		schema.Comment("系统游戏分类表"),                 // 设置数据库表注释
		entsql.Annotation{Table: "game_category"}, // 设置数据库表名
	}
}
