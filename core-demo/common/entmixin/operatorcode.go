package entmixin

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// OperatorCodeMixin 分站编码混合
type OperatorCodeMixin struct{ mixin.Schema }

// Fields 分站编码字段
func (OperatorCodeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("operator_code").MaxLen(64).Optional().Nillable().Comment("Operator Code | 分站编码"),
	}
}

// Interceptors 分站编码拦截器
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

// Hooks 分站编码钩子
func (OperatorCodeMixin) Hooks() []ent.Hook {
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
					if c.OperatorCode != "" {
						if _, exists := m.Field("operator_code"); !exists {
							if err := m.SetField("operator_code", c.OperatorCode); err != nil {
								return nil, err
							}
						}
					}
					return next.Mutate(ctx, m)
				}
				if w, ok := m.(interface {
					WhereP(...func(*sql.Selector))
				}); ok {
					w.WhereP(operatorCodeWhere(c.OperatorCode))
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}

// FilterOperatorCode 过滤分站编码
func FilterOperatorCode(ctx context.Context, q Query) {
	if ctxdata.TenantSkipped(ctx) {
		return
	}
	c := ctxdata.ClaimsFromCtx(ctx)
	if c == nil {
		return
	}
	q.WhereP(operatorCodeWhere(c.OperatorCode))
}

func operatorCodeWhere(code string) func(*sql.Selector) {
	return func(s *sql.Selector) {
		if code == "" {
			s.Where(sql.IsNull(s.C("operator_code")))
			return
		}
		s.Where(sql.EQ(s.C("operator_code"), code))
	}
}
