package entmixin

import "entgo.io/ent/dialect/sql"

// Query is the generated intercept.Query surface (WhereP).
// *UserQuery / *PromotionQuery do not implement WhereP; only the intercept wrapper does.
type Query interface {
	WhereP(...func(*sql.Selector))
}
