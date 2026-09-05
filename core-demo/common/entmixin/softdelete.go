package entmixin

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type SoftDeleteMixin struct{ mixin.Schema }

func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			w, ok := q.(Query)
			if !ok {
				return nil
			}
			FilterSoftDelete(ctx, w)
			return nil
		}),
	}
}

func (SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
					return next.Mutate(ctx, m)
				}
				if ctxdata.SoftDeleteSkipped(ctx) {
					return next.Mutate(ctx, m)
				}
				mx, ok := m.(interface {
					SetOp(ent.Op)
					SetDeletedAt(time.Time)
					WhereP(...func(*sql.Selector))
				})
				if !ok {
					return nil, fmt.Errorf("soft delete: unexpected mutation %T", m)
				}
				mx.WhereP(func(s *sql.Selector) {
					s.Where(sql.IsNull(s.C("deleted_at")))
				})
				mx.SetOp(ent.OpUpdate)
				mx.SetDeletedAt(time.Now())
				return next.Mutate(ctx, m)
			})
		},
	}
}

func FilterSoftDelete(ctx context.Context, q Query) {
	if ctxdata.SoftDeleteSkipped(ctx) {
		return
	}
	q.WhereP(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C("deleted_at")))
	})
}
