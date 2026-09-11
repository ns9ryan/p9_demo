package ent

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GameModel = (*customGameModel)(nil)

type (
	// GameModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameModel.
	GameModel interface {
		gameModel
		withSession(session sqlx.Session) GameModel
	}

	customGameModel struct {
		*defaultGameModel
	}
)

// NewGameModel returns a model for the database table.
func NewGameModel(conn sqlx.SqlConn) GameModel {
	return &customGameModel{
		defaultGameModel: newGameModel(conn),
	}
}

func (m *customGameModel) withSession(session sqlx.Session) GameModel {
	return NewGameModel(sqlx.NewSqlConnFromSession(session))
}

// TableName specifies the table name for GORM
func (Game) TableName() string {
	return "game"
}
