package sync

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// CheckpointManager 检查点管理器
type CheckpointManager struct {
	db *gorm.DB
}

// NewCheckpointManager 创建检查点管理器
func NewCheckpointManager(db *gorm.DB) *CheckpointManager {
	return &CheckpointManager{db: db}
}

// buildCheckpoint 构建检查点对象（从 UpdateCheckpoint 提取的数据构建逻辑）
func (m *CheckpointManager) buildCheckpoint(
	syncScope string,
	checkpointValue string,
	syncResult *vendors.SyncRunResp,
	syncErr error,
) *model.GameSyncCheckpoint {
	now := time.Now()

	// 如果 checkpointValue 为空，使用当前时间戳作为默认值
	if checkpointValue == "" {
		checkpointValue = now.Format(time.RFC3339Nano)
		logger.Debugf("[检查点管理] checkpoint_value 为空，使用当前时间戳: %s", checkpointValue)
	}

	checkpoint := &model.GameSyncCheckpoint{
		SyncScope:       syncScope,
		CheckpointValue: checkpointValue,
		LastSyncAt:      now,
	}

	// 处理同步状态
	if syncErr != nil {
		checkpoint.SyncStatus = model.SyncStatusFailed
		errMsg := syncErr.Error()
		checkpoint.LastErrorMessage = &errMsg
		logger.Errorf("[检查点管理] 同步失败，保存错误信息: %s", syncErr.Error())
	} else if syncResult == nil {
		checkpoint.SyncStatus = model.SyncStatusFailed
		errMsg := "sync result is nil"
		checkpoint.LastErrorMessage = &errMsg
		logger.Error("[检查点管理] 同步结果为空")
	} else {
		checkpoint.SyncStatus = model.SyncStatusSuccess
		checkpoint.LastSuccessAt = now

		// 记录同步统计信息
		if syncResult.Preview != nil && syncResult.Preview.Stats != nil {
			checkpoint.RemoteTotal = syncResult.Preview.Stats.RemoteTotal
			checkpoint.LocalTotal = syncResult.Preview.Stats.LocalTotal
		}

		if syncResult.Apply != nil {
			checkpoint.CreatedCount = int64(syncResult.Apply.Created)
			checkpoint.UpdatedCount = int64(syncResult.Apply.Updated)
			checkpoint.DeletedCount = int64(syncResult.Apply.Deleted)
			checkpoint.FailedCount = int64(syncResult.Apply.Failed)
		}

		logger.Infof("[检查点管理] 同步成功 - 新增:%d, 更新:%d, 删除:%d, 失败:%d",
			checkpoint.CreatedCount, checkpoint.UpdatedCount, checkpoint.DeletedCount, checkpoint.FailedCount)
	}

	return checkpoint
}

// saveOrUpdateCheckpoint 保存检查点到数据库（只做新增操作）
func (m *CheckpointManager) saveOrUpdateCheckpoint(ctx context.Context, checkpoint *model.GameSyncCheckpoint) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 每次同步都插入新的检查点记录，不做更新操作
		if err := tx.Create(checkpoint).Error; err != nil {
			logger.Errorf("[检查点管理] 创建检查点失败: %v", err)
			return err
		}
		logger.Infof("[检查点管理] ✓ 检查点创建成功: scope=%s, id=%d", checkpoint.SyncScope, checkpoint.ID)
		return nil
	})
}

// UpdateCheckpoint 更新同步检查点
func (m *CheckpointManager) UpdateCheckpoint(
	ctx context.Context,
	syncScope string,
	checkpointValue string,
	syncResult *vendors.SyncRunResp,
	syncErr error,
) error {
	logger.Infof("[检查点管理] 准备更新检查点: scope=%s, value=%s", syncScope, checkpointValue)

	// 构建检查点对象
	checkpoint := m.buildCheckpoint(syncScope, checkpointValue, syncResult, syncErr)

	// 保存或更新检查点
	return m.saveOrUpdateCheckpoint(ctx, checkpoint)
}

// GetCheckpoint 获取同步检查点（返回最新的那条记录）
func (m *CheckpointManager) GetCheckpoint(ctx context.Context, syncScope string) (*model.GameSyncCheckpoint, error) {
	var checkpoint model.GameSyncCheckpoint
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

// GetAllCheckpoints 获取所有同步检查点
func (m *CheckpointManager) GetAllCheckpoints(ctx context.Context) ([]*model.GameSyncCheckpoint, error) {
	var checkpoints []*model.GameSyncCheckpoint
	err := m.db.WithContext(ctx).Order("updated_at DESC").Find(&checkpoints).Error
	if err != nil {
		logger.Errorf("[检查点管理] 获取所有检查点失败: %v", err)
		return nil, err
	}
	return checkpoints, nil
}

// DeleteCheckpoint 删除同步检查点
func (m *CheckpointManager) DeleteCheckpoint(ctx context.Context, syncScope string) error {
	err := m.db.WithContext(ctx).Where("sync_scope = ?", syncScope).Delete(&model.GameSyncCheckpoint{}).Error
	if err != nil {
		logger.Errorf("[检查点管理] 删除检查点失败: %v", err)
		return err
	}
	logger.Infof("[检查点管理] ✓ 检查点删除成功: scope=%s", syncScope)
	return nil
}
