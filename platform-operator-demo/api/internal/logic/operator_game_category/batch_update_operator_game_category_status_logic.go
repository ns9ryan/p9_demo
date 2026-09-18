// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_category

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateOperatorGameCategoryStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量修改分站游戏分类状态
func NewBatchUpdateOperatorGameCategoryStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameCategoryStatusLogic {
	return &BatchUpdateOperatorGameCategoryStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateOperatorGameCategoryStatusLogic) BatchUpdateOperatorGameCategoryStatus(req *types.BatchUpdateOperatorGameCategoryStatusRequest) (resp *types.BatchUpdateOperatorGameCategoryStatusResponse, err error) {
	// 调用 RPC 批量修改分站游戏分类状态
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameCategoryServiceClient().BatchUpdateOperatorGameCategoryStatus(l.ctx, &platform_game.BatchUpdateOperatorGameCategoryStatusRequest{
		Ids:    req.Ids,
		Status: req.Status,
	})
	if err != nil {
		l.Logger.Error("BatchUpdateOperatorGameCategoryStatus error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchUpdateOperatorGameCategoryStatusResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
