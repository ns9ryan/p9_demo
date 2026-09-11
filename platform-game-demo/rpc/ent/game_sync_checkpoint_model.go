package ent

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GameSyncCheckpointModel = (*customGameSyncCheckpointModel)(nil)

type (
	// GameSyncCheckpointModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGameSyncCheckpointModel.
	GameSyncCheckpointModel interface {
		gameSyncCheckpointModel
		withSession(session sqlx.Session) GameSyncCheckpointModel
	}

	customGameSyncCheckpointModel struct {
		*defaultGameSyncCheckpointModel
	}
)

// NewGameSyncCheckpointModel returns a model for the database table.
func NewGameSyncCheckpointModel(conn sqlx.SqlConn) GameSyncCheckpointModel {
	return &customGameSyncCheckpointModel{
		defaultGameSyncCheckpointModel: newGameSyncCheckpointModel(conn),
	}
}

func (m *customGameSyncCheckpointModel) withSession(session sqlx.Session) GameSyncCheckpointModel {
	return NewGameSyncCheckpointModel(sqlx.NewSqlConnFromSession(session))
}

// TableName specifies the table name for GORM
func (GameSyncCheckpoint) TableName() string {
	return "game_sync_checkpoint"
}
