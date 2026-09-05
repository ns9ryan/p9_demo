package entmixin

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// OperatorCodeMixin 操作员代码混合
type OperatorCodeMixin struct{ mixin.Schema }

// Fields 操作员代码字段
func (OperatorCodeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("operator_code").MaxLen(64).Optional().Nillable(),
	}
}

// Interceptors 操作员代码拦截器
func (OperatorCodeMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			w, ok := q.(Query)
			if !ok {
				return nil
			}
			FilterOperatorCode(ctx, w)
			return nil
		}),
	}
}

// Hooks 操作员代码钩子
func (OperatorCodeMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if ctxdata.TenantSkipped(ctx) {
					return next.Mutate(ctx, m)
				}
				c := ctxdata.ClaimsFromCtx(ctx)
				if c == nil || c.OperatorCode == "" {
					return next.Mutate(ctx, m)
				}
				if m.Op().Is(ent.OpCreate) {
					if _, exists := m.Field("operator_code"); !exists {
						if err := m.SetField("operator_code", c.OperatorCode); err != nil {
							return nil, err
						}
					}
					return next.Mutate(ctx, m)
				}
				if w, ok := m.(interface {
					WhereP(...func(*sql.Selector))
				}); ok {
					w.WhereP(func(s *sql.Selector) {
						s.Where(sql.EQ(s.C("operator_code"), c.OperatorCode))
					})
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}

// FilterOperatorCode 过滤操作员代码
func FilterOperatorCode(ctx context.Context, q Query) {
	if ctxdata.TenantSkipped(ctx) {
		return
	}
	c := ctxdata.ClaimsFromCtx(ctx)
	if c == nil || c.OperatorCode == "" {
		return
	}
	q.WhereP(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C("operator_code"), c.OperatorCode))
	})
}
