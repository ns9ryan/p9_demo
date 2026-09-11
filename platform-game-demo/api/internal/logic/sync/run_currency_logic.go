// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package sync

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunCurrencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRunCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunCurrencyLogic {
	return &RunCurrencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunCurrencyLogic) RunCurrency(req *types.SyncRunReq) (resp *types.SyncRunResp, err error) {
	l.Infof("[API RunCurrency] received sync run request")

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API RunCurrency] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	// 检查频率限制（5秒）
	if !l.svcCtx.SyncRateLimiter.CheckRateLimit("currency") {
		l.Infof("[API RunCurrency] sync request rejected due to rate limit (5s)")
		return nil, fmt.Errorf("sync operation is already in progress, please wait 5 seconds before retry")
	}

	// 调用 RPC 的 SyncRun，由 RPC 侧负责创建 checkpoint 和异步处理
	client := l.svcCtx.GrpcClient.GetSyncServiceClient()
	runReq := &platformgame.SyncRunRequest{
		ObjectType: "currency",
		SyncCols:   req.SyncCols,
	}

	runResp, err := client.SyncRun(l.ctx, runReq)
	if err != nil {
		l.Errorf("[API RunCurrency] SyncRun call failed: %v", err)
		return nil, fmt.Errorf("sync run failed: %v", err)
	}

	if runResp == nil {
		l.Error("[API RunCurrency] SyncRun response is nil")
		return nil, fmt.Errorf("sync run response is nil")
	}

	l.Infof("[API RunCurrency] sync run success: checkpoint_id=%d", runResp.CheckpointId)

	// 返回 checkpoint ID 给前端
	resp = &types.SyncRunResp{
		CheckpointID: runResp.CheckpointId,
	}

	return resp, nil
}
