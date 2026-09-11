package ent

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GameCategoryModel = (*customGameCategoryModel)(nil)

type (
	// GameCategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameCategoryModel.
	GameCategoryModel interface {
		gameCategoryModel
		withSession(session sqlx.Session) GameCategoryModel
	}

	customGameCategoryModel struct {
		*defaultGameCategoryModel
	}
)

// NewGameCategoryModel returns a model for the database table.
func NewGameCategoryModel(conn sqlx.SqlConn) GameCategoryModel {
	return &customGameCategoryModel{
		defaultGameCategoryModel: newGameCategoryModel(conn),
	}
}

func (m *customGameCategoryModel) withSession(session sqlx.Session) GameCategoryModel {
	return NewGameCategoryModel(sqlx.NewSqlConnFromSession(session))
}

// TableName specifies the table name for GORM
func (GameCategory) TableName() string {
	return "game_category"
}
