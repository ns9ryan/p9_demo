package sync

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// SyncService 游戏数据同步服务接口
type SyncService interface {
	// Preview 预检查（不修改数据，仅显示差异）
	Preview(ctx context.Context, objectType string) (*vendors.SyncPreviewResp, error)

	// Run 执行同步（可选自动应用到数据库）
	Run(ctx context.Context, objectType string, autoApply bool) (*vendors.SyncRunResp, error)

	// SyncAll 全量同步所有对象类型
	SyncAll(ctx context.Context, autoApply bool) (*vendors.SyncRunResp, error)
}
