package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// SortMixin 通用排序字段
type SortMixin struct {
	mixin.Schema
}

// Fields 定义通用排序字段
func (SortMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("sort_no").
			SchemaType(map[string]string{
				dialect.Postgres: "integer",
			}).
			Comment("排序值, 数值越小越靠前"),
	}
}
