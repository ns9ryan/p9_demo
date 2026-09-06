package model

import "time"

// GameSyncCheckpoint 游戏同步进度检查点
type GameSyncCheckpoint struct {
	ID               int64     `gorm:"primaryKey;column:id" json:"id"`
	SyncScope        string    `gorm:"column:sync_scope;type:varchar(64)" json:"sync_scope"`
	CheckpointValue  string    `gorm:"column:checkpoint_value;type:varchar(255)" json:"checkpoint_value"`
	RemoteTotal      int64     `gorm:"column:remote_total;default:0" json:"remote_total"`
	LocalTotal       int64     `gorm:"column:local_total;default:0" json:"local_total"`
	CreatedCount     int64     `gorm:"column:created_count;default:0" json:"created_count"`
	UpdatedCount     int64     `gorm:"column:updated_count;default:0" json:"updated_count"`
	DeletedCount     int64     `gorm:"column:deleted_count;default:0" json:"deleted_count"`
	FailedCount      int64     `gorm:"column:failed_count;default:0" json:"failed_count"`
	SyncStatus       int16     `gorm:"column:sync_status;default:0" json:"sync_status"` // 0: failed, 1: success
	LastErrorMessage *string   `gorm:"column:last_error_message;type:text" json:"last_error_message"`
	LastSuccessAt    time.Time `gorm:"column:last_success_at;type:timestamptz(3);default:CURRENT_TIMESTAMP" json:"last_success_at"`
	LastSyncAt       time.Time `gorm:"column:last_sync_at;type:timestamptz(3);default:CURRENT_TIMESTAMP" json:"last_sync_at"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamptz(3);default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamptz(3);default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 设置表名
func (GameSyncCheckpoint) TableName() string {
	return "game_sync_checkpoint"
}

// SyncScope 常量定义
const (
	SyncScopeCategory     = "CATEGORY"
	SyncScopeProvider     = "PROVIDER"
	SyncScopeChannel      = "CHANNEL"
	SyncScopeGame         = "GAME"
	SyncScopeGameCurrency = "GAME_CURRENCY"
)

// SyncStatus 常量定义
const (
	SyncStatusFailed  int16 = 0
	SyncStatusSuccess int16 = 1
)
