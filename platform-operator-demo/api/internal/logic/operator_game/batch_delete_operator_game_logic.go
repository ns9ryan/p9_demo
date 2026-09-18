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

type BatchDeleteOperatorGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除分站游戏
func NewBatchDeleteOperatorGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperatorGameLogic {
	return &BatchDeleteOperatorGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteOperatorGameLogic) BatchDeleteOperatorGame(req *types.BatchDeleteOperatorGameRequest) (resp *types.BatchDeleteOperatorGameResponse, err error) {
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameServiceClient().BatchDeleteOperatorGame(l.ctx, &platform_game.BatchDeleteOperatorGameRequest{
		Ids: req.Ids,
	})
	if err != nil {
		l.Logger.Error("BatchDeleteOperatorGame error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchDeleteOperatorGameResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
