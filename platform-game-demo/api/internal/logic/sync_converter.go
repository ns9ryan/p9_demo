// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"time"

	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// convertToIsDeleted 将 DeletedAt 转换为 IsDeleted
// nil 表示未删除（0），非nil 表示已删除（1）
func convertToIsDeleted(deletedAt *time.Time) int16 {
	if deletedAt == nil {
		return 0
	}
	return 1
}

// convertSyncPreviewResp 将 vendors.SyncPreviewResp 转换为 types.SyncPreviewResp
func convertSyncPreviewResp(vendorsResp *vendors.SyncPreviewResp) *types.SyncPreviewResp {
	if vendorsResp == nil {
		return nil
	}

	resp := &types.SyncPreviewResp{
		Stats: types.SyncStats{
			RemoteTotal:   vendorsResp.Stats.RemoteTotal,
			LocalTotal:    vendorsResp.Stats.LocalTotal,
			CreateTotal:   vendorsResp.Stats.CreateTotal,
			UpdateTotal:   vendorsResp.Stats.UpdateTotal,
			DeleteTotal:   vendorsResp.Stats.DeleteTotal,
			NoopTotal:     vendorsResp.Stats.NoopTotal,
			ConflictTotal: vendorsResp.Stats.ConflictTotal,
			ErrorTotal:    vendorsResp.Stats.ErrorTotal,
		},
		Diffs: make([]types.SyncDiff, len(vendorsResp.Diffs)),
	}

	for i, diff := range vendorsResp.Diffs {
		resp.Diffs[i] = types.SyncDiff{
			ObjectType:   diff.ObjectType,
			ObjectID:     diff.ObjectId,
			ObjectCode:   diff.ObjectCode,
			Action:       diff.Action,
			RemoteID:     diff.RemoteId,
			RemoteCode:   diff.RemoteCode,
			Reason:       diff.Reason,
			ConflictType: diff.ConflictType,
			Details:      diff.Details,
		}
	}

	return resp
}

// convertSyncRunResp 将 vendors.SyncRunResp 转换为 types.SyncRunResp
func convertSyncRunResp(vendorsResp *vendors.SyncRunResp) *types.SyncRunResp {
	if vendorsResp == nil {
		return nil
	}

	previewResp := convertSyncPreviewResp(vendorsResp.Preview)
	if previewResp == nil {
		previewResp = &types.SyncPreviewResp{
			Stats: types.SyncStats{},
			Diffs: make([]types.SyncDiff, 0),
		}
	}

	resp := &types.SyncRunResp{
		Preview: *previewResp,
		Apply: types.SyncApplyResult{
			Created: vendorsResp.Apply.Created,
			Updated: vendorsResp.Apply.Updated,
			Deleted: vendorsResp.Apply.Deleted,
			Failed:  vendorsResp.Apply.Failed,
			Skipped: vendorsResp.Apply.Skipped,
		},
	}

	return resp
}
