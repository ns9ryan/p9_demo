// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateOperatorGameStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量修改分站游戏状态
func NewBatchUpdateOperatorGameStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameStatusLogic {
	return &BatchUpdateOperatorGameStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateOperatorGameStatusLogic) BatchUpdateOperatorGameStatus(req *types.BatchUpdateOperatorGameStatusRequest) (resp *types.BatchUpdateOperatorGameStatusResponse, err error) {
	// 调用 RPC 批量修改分站游戏状态
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameServiceClient().BatchUpdateOperatorGameStatus(l.ctx, &platform_game.BatchUpdateOperatorGameStatusRequest{
		Ids:    req.Ids,
		Status: req.Status,
	})
	if err != nil {
		l.Logger.Error("BatchUpdateOperatorGameStatus error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchUpdateOperatorGameStatusResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
