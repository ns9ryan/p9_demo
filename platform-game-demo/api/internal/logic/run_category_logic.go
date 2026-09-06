// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRunCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunCategoryLogic {
	return &RunCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunCategoryLogic) RunCategory(req *types.SyncRunReq) (resp *types.SyncRunResp, err error) {
	l.Infof("[API RunCategory] received sync run request")

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API RunCategory] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	client := l.svcCtx.GrpcClient.GetPlatformGameServiceClient()

	// 第一步：调用 SyncPreview 获取预检查结果（包含 diffs）
	previewReq := &platformgame.SyncPreviewRequest{
		ObjectType: "category",
	}
	previewResp, err := client.SyncPreview(l.ctx, previewReq)
	if err != nil {
		l.Errorf("[API RunCategory] SyncPreview call failed: %v", err)
		return nil, err
	}

	if previewResp == nil {
		l.Error("[API RunCategory] SyncPreview response is nil")
		return nil, fmt.Errorf("SyncPreview response is nil")
	}

	if previewResp.Code != constant.CodeSuccess {
		l.Errorf("[API RunCategory] SyncPreview error: Code=%d, Message=%s", previewResp.Code, previewResp.Message)
		return nil, fmt.Errorf("SyncPreview error: %s", previewResp.Message)
	}

	// 第二步：调用 SyncRun 执行同步
	runReq := &platformgame.SyncRunRequest{
		ObjectType: "category",
		AutoApply:  true,
	}
	runResp, err := client.SyncRun(l.ctx, runReq)
	if err != nil {
		l.Errorf("[API RunCategory] SyncRun call failed: %v", err)
		return nil, err
	}

	if runResp == nil {
		l.Error("[API RunCategory] SyncRun response is nil")
		return nil, fmt.Errorf("SyncRun response is nil")
	}

	if runResp.Code != constant.CodeSuccess {
		l.Errorf("[API RunCategory] SyncRun error: Code=%d, Message=%s", runResp.Code, runResp.Message)
		return nil, fmt.Errorf("SyncRun error: %s", runResp.Message)
	}

	// 构建 API 响应，包含 Preview 和 Apply 结果
	previewResp2 := types.SyncPreviewResp{
		Stats: types.SyncStats{
			RemoteTotal: int64(previewResp.Stats.Total),
			LocalTotal:  int64(previewResp.Stats.Total),
			CreateTotal: int64(previewResp.Stats.Added),
			UpdateTotal: int64(previewResp.Stats.Updated),
			DeleteTotal: int64(previewResp.Stats.Deleted),
		},
	}

	// 转换 diffs
	if previewResp.Diffs != nil && len(previewResp.Diffs) > 0 {
		previewResp2.Diffs = make([]types.SyncDiff, 0, len(previewResp.Diffs))
		for _, diff := range previewResp.Diffs {
			previewResp2.Diffs = append(previewResp2.Diffs, types.SyncDiff{
				ObjectType: diff.Type,
				Action:     diff.Action,
				ObjectCode: diff.ObjectId,
				Details:    diff.Description,
			})
		}
	} else {
		previewResp2.Diffs = make([]types.SyncDiff, 0)
	}

	resp = &types.SyncRunResp{
		Preview: previewResp2,
		Apply: types.SyncApplyResult{
			Created: runResp.Created,
			Updated: runResp.Updated,
			Deleted: runResp.Deleted,
			Failed:  runResp.Failed,
			Skipped: runResp.Skipped,
		},
	}

	l.Infof("[API RunCategory] success: preview_total=%d, preview_diffs=%d, created=%d, updated=%d, deleted=%d, failed=%d, skipped=%d",
		previewResp.Stats.Total, len(previewResp2.Diffs), runResp.Created, runResp.Updated, runResp.Deleted, runResp.Failed, runResp.Skipped)
	return resp, nil
}
