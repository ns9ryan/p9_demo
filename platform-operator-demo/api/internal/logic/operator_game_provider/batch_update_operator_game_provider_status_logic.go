// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_provider

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateOperatorGameProviderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量修改分站游戏提供商状态
func NewBatchUpdateOperatorGameProviderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameProviderStatusLogic {
	return &BatchUpdateOperatorGameProviderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateOperatorGameProviderStatusLogic) BatchUpdateOperatorGameProviderStatus(req *types.BatchUpdateOperatorGameProviderStatusRequest) (resp *types.BatchUpdateOperatorGameProviderStatusResponse, err error) {
	// 调用 RPC 批量修改分站游戏提供商状态
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameProviderServiceClient().BatchUpdateOperatorGameProviderStatus(l.ctx, &platform_game.BatchUpdateOperatorGameProviderStatusRequest{
		Ids:    req.Ids,
		Status: req.Status,
	})
	if err != nil {
		l.Logger.Error("BatchUpdateOperatorGameProviderStatus error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchUpdateOperatorGameProviderStatusResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
