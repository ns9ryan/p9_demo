package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type CasbinRule struct{ ent.Schema }

func (CasbinRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("Casbin Rule Table | 权限策略表"),
		entsql.Annotation{Table: "casbin_rule"},
	}
}

func (CasbinRule) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Primary key | 主键"),
		field.String("ptype").MaxLen(16).Default("p").Comment("Policy type | 策略类型"),
		field.String("v0").MaxLen(255).Default("").Comment("Subject / role | 主体 / 角色"),
		field.String("v1").MaxLen(255).Default("").Comment("Domain | 域"),
		field.String("v2").MaxLen(255).Default("").Comment("Object / path | 对象 / 路径"),
		field.String("v3").MaxLen(255).Default("").Comment("Action / method | 操作 / 方法"),
		field.String("v4").MaxLen(255).Default("").Comment("Extra | 扩展"),
		field.String("v5").MaxLen(255).Default("").Comment("Extra | 扩展"),
	}
}

func (CasbinRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ptype", "v1", "v0", "v2", "v3", "v4", "v5").Unique().
			StorageKey("uk_casbin_rule_policy"),
	}
}
