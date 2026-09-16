package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// IDMixin 通用主键字段
type IDMixin struct {
	mixin.Schema
}

// Fields 定义通用主键字段
func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			SchemaType(map[string]string{
				dialect.Postgres: "bigint",
			}).
			Comment("主键ID"),
	}
}
