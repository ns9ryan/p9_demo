package sync

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

// SyncService 游戏数据同步服务接口
type SyncService interface {
	// Preview 预检查（不修改数据，仅显示差异）
	Preview(ctx context.Context, objectType string) (*platform_game.SyncPreviewResp, error)

	// Run 执行同步（可选自动应用到数据库）
	Run(ctx context.Context, objectType string, autoApply bool) (*platform_game.SyncRunResp, error)
}
