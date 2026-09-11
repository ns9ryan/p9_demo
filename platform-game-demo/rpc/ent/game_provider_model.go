package ent

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GameProviderModel = (*customGameProviderModel)(nil)

type (
	// GameProviderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameProviderModel.
	GameProviderModel interface {
		gameProviderModel
		withSession(session sqlx.Session) GameProviderModel
	}

	customGameProviderModel struct {
		*defaultGameProviderModel
	}
)

// NewGameProviderModel returns a model for the database table.
func NewGameProviderModel(conn sqlx.SqlConn) GameProviderModel {
	return &customGameProviderModel{
		defaultGameProviderModel: newGameProviderModel(conn),
	}
}

func (m *customGameProviderModel) withSession(session sqlx.Session) GameProviderModel {
	return NewGameProviderModel(sqlx.NewSqlConnFromSession(session))
}

// TableName specifies the table name for GORM
func (GameProvider) TableName() string {
	return "game_provider"
}
