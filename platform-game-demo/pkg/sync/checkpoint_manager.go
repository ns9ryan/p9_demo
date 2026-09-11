package sync

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/rpc/ent"
)

// CheckpointManager 检查点管理器
type CheckpointManager struct {
	db *gorm.DB
}

// NewCheckpointManager 创建检查点管理器
func NewCheckpointManager(db *gorm.DB) *CheckpointManager {
	return &CheckpointManager{db: db}
}

// saveOrUpdateCheckpoint 保存检查点到数据库（只做新增操作）
func (m *CheckpointManager) saveOrUpdateCheckpoint(ctx context.Context, checkpoint *ent.GameSyncCheckpoint) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 每次同步都插入新的检查点记录，不做更新操作
		if err := tx.Create(checkpoint).Error; err != nil {
			logger.Errorf("[检查点管理] 创建检查点失败: %v", err)
			return err
		}
		logger.Infof("[检查点管理] ✓ 检查点创建成功: scope=%s, id=%d", checkpoint.SyncScope, checkpoint.Id)
		return nil
	})
}

// GetCheckpoint 获取同步检查点（返回最新的那条记录）
func (m *CheckpointManager) GetCheckpoint(ctx context.Context, syncScope string) (*ent.GameSyncCheckpoint, error) {
	var checkpoint ent.GameSyncCheckpoint
	err := m.db.WithContext(ctx).
		Where("sync_scope = ?", syncScope).
		Order("created_at DESC, id DESC").
		First(&checkpoint).Error
	if err == gorm.ErrRecordNotFound {
		logger.Debugf("[检查点管理] 检查点不存在: scope=%s", syncScope)
		return nil, nil
	}
	if err != nil {
		logger.Errorf("[检查点管理] 获取检查点失败: %v", err)
		return nil, err
	}
	return &checkpoint, nil
}

// UpdateCheckpointProgress 更新指定检查点的进度
func (m *CheckpointManager) UpdateCheckpointProgress(ctx context.Context, msg *ProgressMessage) error {
	if msg.CheckpointID <= 0 {
		logger.Warnf("[检查点管理] 无效的 checkpointID: %d", msg.CheckpointID)
		return nil
	}

	checkpoint := &ent.GameSyncCheckpoint{
		Progress:        int64(msg.Progress),
		CheckpointValue: fmt.Sprintf("%d", msg.ProcessedCount),
		RemoteTotal:     int64(msg.RemoteTotal),
		LocalTotal:      int64(msg.LocalTotal),
		CreatedCount:    int64(msg.Created),
		UpdatedCount:    int64(msg.Updated),
		DeletedCount:    int64(msg.Deleted),
		FailedCount:     int64(msg.Failed),
		SyncStatus:      1,
		LastSuccessAt:   time.Now(),
		LastSyncAt:      time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := m.db.WithContext(ctx).Model(&ent.GameSyncCheckpoint{}).Where("id = ?", msg.CheckpointID).Updates(checkpoint).Error; err != nil {
		logger.Errorf("[检查点管理] 更新检查点进度失败: checkpointID=%d, progress=%d, error=%v", msg.CheckpointID, msg.Progress, err)
		return err
	}

	logger.Infof("[检查点管理] ✓ 检查点进度已更新: checkpointID=%d, progress=%d%%, processedCount=%d", msg.CheckpointID, msg.Progress, msg.ProcessedCount)
	return nil
}
