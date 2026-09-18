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

type BatchCreateOperatorGameCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建分站游戏分类
func NewBatchCreateOperatorGameCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameCategoryLogic {
	return &BatchCreateOperatorGameCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCreateOperatorGameCategoryLogic) BatchCreateOperatorGameCategory(req *types.BatchCreateOperatorGameCategoryRequest) (resp *types.BatchCreateOperatorGameCategoryResponse, err error) {
	// 构建 RPC 请求
	var items []*platform_game.BatchCreateOperatorGameCategoryItemRequest
	for _, item := range req.Items {
		items = append(items, &platform_game.BatchCreateOperatorGameCategoryItemRequest{
			OpCode:       item.OpCode,
			CategoryCode: item.CategoryCode,
			Status:       item.Status,
		})
	}

	// 调用 RPC 批量创建分站游戏分类
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameCategoryServiceClient().BatchCreateOperatorGameCategory(l.ctx, &platform_game.BatchCreateOperatorGameCategoryRequest{
		Items: items,
	})
	if err != nil {
		l.Logger.Error("BatchCreateOperatorGameCategory error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchCreateOperatorGameCategoryResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
