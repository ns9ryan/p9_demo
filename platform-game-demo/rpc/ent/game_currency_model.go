package ent

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GameCurrencyModel = (*customGameCurrencyModel)(nil)

type (
	// GameCurrencyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameCurrencyModel.
	GameCurrencyModel interface {
		gameCurrencyModel
		withSession(session sqlx.Session) GameCurrencyModel
	}

	customGameCurrencyModel struct {
		*defaultGameCurrencyModel
	}
)

// NewGameCurrencyModel returns a model for the database table.
func NewGameCurrencyModel(conn sqlx.SqlConn) GameCurrencyModel {
	return &customGameCurrencyModel{
		defaultGameCurrencyModel: newGameCurrencyModel(conn),
	}
}

func (m *customGameCurrencyModel) withSession(session sqlx.Session) GameCurrencyModel {
	return NewGameCurrencyModel(sqlx.NewSqlConnFromSession(session))
}

// TableName specifies the table name for GORM
func (GameCurrency) TableName() string {
	return "game_currency"
}
