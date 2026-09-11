package ent

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GameChannelModel = (*customGameChannelModel)(nil)

type (
	// GameChannelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameChannelModel.
	GameChannelModel interface {
		gameChannelModel
		withSession(session sqlx.Session) GameChannelModel
	}

	customGameChannelModel struct {
		*defaultGameChannelModel
	}
)

// NewGameChannelModel returns a model for the database table.
func NewGameChannelModel(conn sqlx.SqlConn) GameChannelModel {
	return &customGameChannelModel{
		defaultGameChannelModel: newGameChannelModel(conn),
	}
}

func (m *customGameChannelModel) withSession(session sqlx.Session) GameChannelModel {
	return NewGameChannelModel(sqlx.NewSqlConnFromSession(session))
}

// TableName specifies the table name for GORM
func (GameChannel) TableName() string {
	return "game_channel"
}
