package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// AllocationStatusMixin 通用资源分配状态字段
type AllocationStatusMixin struct {
	mixin.Schema
}

// Fields 定义通用资源分配状态字段
func (AllocationStatusMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("allocation_status").
			Default(1).
			Range(1, 2).
			SchemaType(map[string]string{
				dialect.Postgres: "smallint",
			}).
			Comment("分配状态: 1已分配, 2已回收"),
	}
}
