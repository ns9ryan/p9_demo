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

type BatchDeleteOperatorGameCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除分站游戏分类
func NewBatchDeleteOperatorGameCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperatorGameCategoryLogic {
	return &BatchDeleteOperatorGameCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteOperatorGameCategoryLogic) BatchDeleteOperatorGameCategory(req *types.BatchDeleteOperatorGameCategoryRequest) (resp *types.BatchDeleteOperatorGameCategoryResponse, err error) {
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameCategoryServiceClient().BatchDeleteOperatorGameCategory(l.ctx, &platform_game.BatchDeleteOperatorGameCategoryRequest{
		Ids: req.Ids,
	})
	if err != nil {
		l.Logger.Error("BatchDeleteOperatorGameCategory error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchDeleteOperatorGameCategoryResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
