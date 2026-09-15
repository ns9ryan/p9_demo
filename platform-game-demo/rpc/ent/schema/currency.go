package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Currency holds the schema definition for the Currency entity.
type Currency struct {
	ent.Schema
}

// Table of the Currency.
func (Currency) Table() string {
	return "currency"
}

// Fields of the Currency.
func (Currency) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.Int64("status"),
		field.Int64("sort_no"),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.String("code"),
		field.String("name_key"),
		field.Int64("currency_type").Range(1, 2),
		field.String("symbol"),
		field.Int64("amount_factor"),
	}
}

// Edges of the Currency.
func (Currency) Edges() []ent.Edge {
	return nil
}

func (Currency) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),            // 启用数据库字段注释
		schema.Comment("系统货币表"),              // 设置数据库表注释
		entsql.Annotation{Table: "currency"}, // 设置数据库表名
	}
}
