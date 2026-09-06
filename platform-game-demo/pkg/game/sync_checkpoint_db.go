package game

import (
	"context"
	"fmt"
	"log"
	"time"

	"oa.98ent.com/p9/platform-game/common/model"

	"gorm.io/gorm"
)

// GameSyncCheckpointGetParams 定义 GameSyncCheckpointGet 的参数
type GameSyncCheckpointGetParams struct {
	ID        int64
	SyncScope string
}

// GameSyncCheckpointGetDB 获取单个同步检查点
func GameSyncCheckpointGetDB(ctx context.Context, db *gorm.DB, params *GameSyncCheckpointGetParams) (*model.GameSyncCheckpoint, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	query := db.WithContext(ctx)

	// 如果指定了 sync_scope，按 sync_scope 查询；否则获取最近一条记录
	if params.SyncScope != "" {
		query = query.Where("sync_scope = ?", params.SyncScope)
	} else if params.ID > 0 {
		query = query.Where("id = ?", params.ID)
	} else {
		query = query.Order("last_sync_at DESC")
	}

	var checkpoint model.GameSyncCheckpoint
	if err := query.First(&checkpoint).Error; err != nil {
		return nil, fmt.Errorf("[GameSyncCheckpointGetDB] query failed: %v", err)
	}

	return &checkpoint, nil
}

// GameSyncCheckpointListParams 定义 GameSyncCheckpointList 的参数
type GameSyncCheckpointListParams struct {
	Page      int64
	PageSize  int64
	SyncScope string
	StartTime int64 // Unix timestamp
	EndTime   int64 // Unix timestamp
}

// GameSyncCheckpointListResult 定义 GameSyncCheckpointList 的返回结果
type GameSyncCheckpointListResult struct {
	Checkpoints []*model.GameSyncCheckpoint
	Total       int64
}

// GameSyncCheckpointListDB 分页列表查询同步检查点
func GameSyncCheckpointListDB(ctx context.Context, db *gorm.DB, params *GameSyncCheckpointListParams) (*GameSyncCheckpointListResult, error) {
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 15
	}
	if pageSize > 100 {
		pageSize = 100
	}
	log.Printf("[GameSyncCheckpointListDB] input params: page=%d, page_size=%d, sync_scope=%q, start_time=%d, end_time=%d", params.Page, params.PageSize, params.SyncScope, params.StartTime, params.EndTime)

	// 构建查询条件
	query := db.WithContext(ctx)

	// 按 sync_scope 筛选（可选）
	if params.SyncScope != "" {
		log.Printf("[GameSyncCheckpointListDB] apply filter: sync_scope = %q", params.SyncScope)
		query = query.Where("sync_scope = ?", params.SyncScope)
	}

	// 按 last_success_at 时间区间筛选（可选）
	if params.StartTime > 0 {
		startAt := time.Unix(params.StartTime, 0)
		log.Printf("[GameSyncCheckpointListDB] apply filter: last_success_at >= %s", startAt.Format(time.RFC3339))
		query = query.Where("last_success_at >= ?", startAt)
	}
	if params.EndTime > 0 {
		endAt := time.Unix(params.EndTime, 0)
		log.Printf("[GameSyncCheckpointListDB] apply filter: last_success_at <= %s", endAt.Format(time.RFC3339))
		query = query.Where("last_success_at <= ?", endAt)
	}

	countDryRun := query.Session(&gorm.Session{DryRun: true}).Model(&model.GameSyncCheckpoint{}).Count(new(int64))
	log.Printf("[GameSyncCheckpointListDB] count SQL: %s | vars=%v", countDryRun.Statement.SQL.String(), countDryRun.Statement.Vars)

	var total int64
	if err := query.
		Model(&model.GameSyncCheckpoint{}).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("[GameSyncCheckpointListDB] count failed: %v", err)
	}
	log.Printf("[GameSyncCheckpointListDB] count result: total=%d", total)

	listDryRun := query.Session(&gorm.Session{DryRun: true}).
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order("last_sync_at DESC").
		Find(&[]*model.GameSyncCheckpoint{})
	log.Printf("[GameSyncCheckpointListDB] list SQL: %s | vars=%v", listDryRun.Statement.SQL.String(), listDryRun.Statement.Vars)

	var checkpoints []*model.GameSyncCheckpoint
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order("last_sync_at DESC").
		Find(&checkpoints).Error; err != nil {
		return nil, fmt.Errorf("[GameSyncCheckpointListDB] query failed: %v", err)
	}
	log.Printf("[GameSyncCheckpointListDB] list result: rows=%d", len(checkpoints))

	return &GameSyncCheckpointListResult{
		Checkpoints: checkpoints,
		Total:       total,
	}, nil
}
