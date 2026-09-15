package dao

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/ent/gamesynccheckpoint"
)

type GameSyncCheckpointDAO struct {
	db *ent.Client
}

func NewGameSyncCheckpointDAO(db *ent.Client) *GameSyncCheckpointDAO {
	return &GameSyncCheckpointDAO{
		db: db,
	}
}

// GetGameSyncCheckpointByID 根据ID获取同步检查点
func (d *GameSyncCheckpointDAO) GetGameSyncCheckpointByID(ctx context.Context, id int64) (*ent.GameSyncCheckpoint, error) {
	return d.db.GameSyncCheckpoint.Query().
		Where(gamesynccheckpoint.IDEQ(id)).
		Only(ctx)
}

// GetGameSyncCheckpointList 获取同步检查点列表
func (d *GameSyncCheckpointDAO) GetGameSyncCheckpointList(ctx context.Context, offset, limit int64) ([]*ent.GameSyncCheckpoint, int, error) {
	query := d.db.GameSyncCheckpoint.Query()

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// 获取分页数据
	checkpoints, err := query.
		Order(gamesynccheckpoint.ByID()).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	return checkpoints, total, nil
}

func (d *GameSyncCheckpointDAO) CreateGameSyncCheckpoint(ctx context.Context, checkpoint *ent.GameSyncCheckpoint) (*ent.GameSyncCheckpoint, error) {
	return d.db.GameSyncCheckpoint.Create().
		SetSyncScope(checkpoint.SyncScope).
		SetSyncStatus(checkpoint.SyncStatus).
		SetCheckpointValue(checkpoint.CheckpointValue).
		SetProgress(checkpoint.Progress).
		SetRemoteTotal(checkpoint.RemoteTotal).
		SetLocalTotal(checkpoint.LocalTotal).
		SetCreatedCount(checkpoint.CreatedCount).
		SetUpdatedCount(checkpoint.UpdatedCount).
		SetDeletedCount(checkpoint.DeletedCount).
		SetFailedCount(checkpoint.FailedCount).
		SetCreatedAt(checkpoint.CreatedAt).
		SetUpdatedAt(checkpoint.UpdatedAt).
		SetLastErrorMessage(checkpoint.LastErrorMessage).
		SetLastSuccessAt(checkpoint.LastSuccessAt).
		SetLastSyncAt(checkpoint.LastSyncAt).
		Save(ctx)
}

func (d *GameSyncCheckpointDAO) UpdateGameSyncCheckpoint(ctx context.Context, updates map[string]interface{}) (*ent.GameSyncCheckpoint, error) {
	update := d.db.GameSyncCheckpoint.UpdateOneID(updates["ID"].(int64))

	// 动态应用更新
	for key, value := range updates {
		switch key {
		case "SyncScope":
			if val, ok := value.(string); ok {
				update = update.SetSyncScope(val)
			}
		case "SyncStatus":
			if val, ok := value.(int64); ok {
				update = update.SetSyncStatus(val)
			}
		case "CheckpointValue":
			if val, ok := value.(string); ok {
				update = update.SetCheckpointValue(val)
			}
		case "Progress":
			if val, ok := value.(int64); ok {
				update = update.SetProgress(val)
			}
		case "RemoteTotal":
			if val, ok := value.(int64); ok {
				update = update.SetRemoteTotal(val)
			}
		case "LocalTotal":
			if val, ok := value.(int64); ok {
				update = update.SetLocalTotal(val)
			}
		case "CreatedCount":
			if val, ok := value.(int64); ok {
				update = update.SetCreatedCount(val)
			}
		case "UpdatedCount":
			if val, ok := value.(int64); ok {
				update = update.SetUpdatedCount(val)
			}
		case "DeletedCount":
			if val, ok := value.(int64); ok {
				update = update.SetDeletedCount(val)
			}
		case "FailedCount":
			if val, ok := value.(int64); ok {
				update = update.SetFailedCount(val)
			}
		case "CreatedAt":
			if val, ok := value.(time.Time); ok {
				update = update.SetCreatedAt(val)
			}
		case "UpdatedAt":
			if val, ok := value.(time.Time); ok {
				update = update.SetUpdatedAt(val)
			}
		}
	}
	return update.Save(ctx)
}

// func (d *GameSyncCheckpointDAO) SaveOrUpdateCheckpoint(ctx context.Context, checkpoint *ent.GameSyncCheckpoint) (*ent.GameSyncCheckpoint, error) {
// 	count, err := d.db.GameSyncCheckpoint.Query().Where(gamesynccheckpoint.IDEQ(checkpoint.ID)).Clone().Count(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if count == 0 {
// 		return d.CreateGameSyncCheckpoint(ctx, checkpoint)
// 	}
// 	updates := map[string]interface{}{
// 		"SyncScope":       checkpoint.SyncScope,
// 		"SyncStatus":      checkpoint.SyncStatus,
// 		"CheckpointValue": checkpoint.CheckpointValue,
// 		"Progress":        checkpoint.Progress,
// 		"RemoteTotal":     checkpoint.RemoteTotal,
// 		"LocalTotal":      checkpoint.LocalTotal,
// 		"CreatedCount":    checkpoint.CreatedCount,
// 		"UpdatedCount":    checkpoint.UpdatedCount,
// 		"DeletedCount":    checkpoint.DeletedCount,
// 		"FailedCount":     checkpoint.FailedCount,
// 		"CreatedAt":       checkpoint.CreatedAt,
// 		"UpdatedAt":       checkpoint.UpdatedAt,
// 	}
// 	return d.UpdateGameSyncCheckpoint(ctx, updates)

// 	// return sql.Insert(gamesynccheckpoint.Table).
// 	// 	Columns(
// 	// 		gamesynccheckpoint.FieldID,
// 	// 		gamesynccheckpoint.FieldSyncScope,
// 	// 		gamesynccheckpoint.FieldSyncStatus,
// 	// 		gamesynccheckpoint.FieldCheckpointValue,
// 	// 		gamesynccheckpoint.FieldProgress,
// 	// 		gamesynccheckpoint.FieldRemoteTotal,
// 	// 		gamesynccheckpoint.FieldLocalTotal,
// 	// 		gamesynccheckpoint.FieldCreatedCount,
// 	// 		gamesynccheckpoint.FieldUpdatedCount,
// 	// 		gamesynccheckpoint.FieldDeletedCount,
// 	// 		gamesynccheckpoint.FieldFailedCount,
// 	// 		gamesynccheckpoint.FieldCreatedAt,
// 	// 		gamesynccheckpoint.FieldUpdatedAt,
// 	// 	).
// 	// 	Values(
// 	// 		checkpoint.ID,
// 	// 		checkpoint.SyncScope,
// 	// 		checkpoint.SyncStatus,
// 	// 		checkpoint.CheckpointValue,
// 	// 		checkpoint.Progress,
// 	// 		checkpoint.RemoteTotal,
// 	// 		checkpoint.LocalTotal,
// 	// 		checkpoint.CreatedCount,
// 	// 		checkpoint.UpdatedCount,
// 	// 		checkpoint.DeletedCount,
// 	// 		checkpoint.FailedCount,
// 	// 		checkpoint.CreatedAt,
// 	// 		checkpoint.UpdatedAt,
// 	// 	).
// 	// 	OnConflict(
// 	// 		sql.ConflictColumns(gamesynccheckpoint.FieldID),
// 	// 	).
// 	// 	Exec(ctx)
// }

func (d *GameSyncCheckpointDAO) FindGameSyncCheckpointBySyncScope(ctx context.Context, syncScope string) (*ent.GameSyncCheckpoint, error) {
	return d.db.GameSyncCheckpoint.Query().
		Where(gamesynccheckpoint.SyncScopeEQ(syncScope)).
		Order(ent.Desc(gamesynccheckpoint.FieldCreatedAt), ent.Desc(gamesynccheckpoint.FieldID)).
		First(ctx)
}
