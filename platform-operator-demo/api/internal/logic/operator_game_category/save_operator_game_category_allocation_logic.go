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

type SaveOperatorGameCategoryAllocationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存分站游戏分类分配
func NewSaveOperatorGameCategoryAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameCategoryAllocationLogic {
	return &SaveOperatorGameCategoryAllocationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveOperatorGameCategoryAllocationLogic) SaveOperatorGameCategoryAllocation(req *types.SaveOperatorGameCategoryAllocationRequest) (resp *types.SaveOperatorGameCategoryAllocationResponse, err error) {
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameCategoryServiceClient().SaveOperatorGameCategoryAllocation(l.ctx, &platform_game.SaveOperatorGameCategoryAllocationRequest{
		OpCode: req.OpCode,
		Items:  convertToRpcSaveOperatorGameCategoryAllocationInfo(req.Items),
	})
	if err != nil {
		return nil, err
	}
	return &types.SaveOperatorGameCategoryAllocationResponse{
		Total:   rpcResp.Total,
		Created: rpcResp.Created,
		Deleted: rpcResp.Deleted,
		Exist:   rpcResp.Exist,
		Failed:  rpcResp.Failed,
	}, nil
}

func convertToRpcSaveOperatorGameCategoryAllocationInfo(items []types.SaveOperatorGameCategoryAllocationItem) []*platform_game.SaveOperatorGameCategoryAllocationInfo {
	rpcItems := make([]*platform_game.SaveOperatorGameCategoryAllocationInfo, len(items))
	for i, item := range items {
		rpcItems[i] = &platform_game.SaveOperatorGameCategoryAllocationInfo{
			Code:        item.Code,
			CheckStatus: item.CheckStatus,
		}
	}
	return rpcItems
}
