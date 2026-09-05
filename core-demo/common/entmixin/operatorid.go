package entmixin

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// OperatorIDMixin 操作员ID混合
type OperatorIDMixin struct{ mixin.Schema }

// Fields 操作员ID字段
func (OperatorIDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("operator_id").Optional().Nillable(),
	}
}

// Interceptors 操作员ID拦截器
func (OperatorIDMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			w, ok := q.(Query)
			if !ok {
				return nil
			}
			FilterOperatorID(ctx, w)
			return nil
		}),
	}
}

// Hooks 操作员ID钩子
func (OperatorIDMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if ctxdata.TenantSkipped(ctx) {
					return next.Mutate(ctx, m)
				}
				c := ctxdata.ClaimsFromCtx(ctx)
				if c == nil {
					return next.Mutate(ctx, m)
				}
				if m.Op().Is(ent.OpCreate) {
					if c.OperatorID != 0 {
						if _, exists := m.Field("operator_id"); !exists {
							if err := m.SetField("operator_id", c.OperatorID); err != nil {
								return nil, err
							}
						}
					}
					return next.Mutate(ctx, m)
				}
				if w, ok := m.(interface {
					WhereP(...func(*sql.Selector))
				}); ok {
					w.WhereP(operatorIDWhere(c.OperatorID))
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}

// FilterOperatorID 过滤操作员ID
func FilterOperatorID(ctx context.Context, q Query) {
	if ctxdata.TenantSkipped(ctx) {
		return
	}
	c := ctxdata.ClaimsFromCtx(ctx)
	if c == nil {
		return
	}
	q.WhereP(operatorIDWhere(c.OperatorID))
}

// operatorIDWhere 操作员ID条件
func operatorIDWhere(id int64) func(*sql.Selector) {
	return func(s *sql.Selector) {
		if id == 0 {
			s.Where(sql.IsNull(s.C("operator_id")))
			return
		}
		s.Where(sql.EQ(s.C("operator_id"), id))
	}
}
