package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Operator holds the schema definition for the Operator entity.
type Operator struct {
	ent.Schema
}

// Table of the Operator.
func (Operator) Table() string {
	return "operator"
}

// Fields of the Operator.
func (Operator) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("code"),
		field.String("name"),
		field.String("timezone_code"),
		field.String("settlement_currency_code"),
		field.Int16("creation_status").Range(1, 2),
		field.Int16("publish_status").Range(1, 2),
		field.String("publish_task_no").Optional(),
		field.Time("published_at").Optional(),
		field.Int16("status").Range(1, 2),
		field.String("remark").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the Operator.
func (Operator) Edges() []ent.Edge {
	return nil
}

func (Operator) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),            // 启用数据库字段注释
		schema.Comment("分站表"),                // 设置数据库表注释
		entsql.Annotation{Table: "operator"}, // 设置数据库表名
	}
}
