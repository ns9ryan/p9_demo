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

type SaveOperatorGameAllocationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存分站游戏分配
func NewSaveOperatorGameAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameAllocationLogic {
	return &SaveOperatorGameAllocationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveOperatorGameAllocationLogic) SaveOperatorGameAllocation(req *types.SaveOperatorGameAllocationRequest) (resp *types.SaveOperatorGameAllocationResponse, err error) {
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameServiceClient().SaveOperatorGameAllocation(l.ctx, &platform_game.SaveOperatorGameAllocationRequest{
		OpCode: req.OpCode,
		Items:  convertToRpcSaveOperatorGameAllocationInfo(req.Items),
	})
	if err != nil {
		return nil, err
	}
	return &types.SaveOperatorGameAllocationResponse{
		Total:   rpcResp.Total,
		Created: rpcResp.Created,
		Deleted: rpcResp.Deleted,
		Exist:   rpcResp.Exist,
		Failed:  rpcResp.Failed,
	}, nil
}

func convertToRpcSaveOperatorGameAllocationInfo(items []types.SaveOperatorGameAllocationItem) []*platform_game.SaveOperatorGameAllocationInfo {
	rpcItems := make([]*platform_game.SaveOperatorGameAllocationInfo, len(items))
	for i, item := range items {
		rpcItems[i] = &platform_game.SaveOperatorGameAllocationInfo{
			Code:        item.Code,
			CheckStatus: item.CheckStatus,
		}
	}
	return rpcItems
}
