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

type SaveOperatorGameProviderAllocationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存分站游戏提供商分配
func NewSaveOperatorGameProviderAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameProviderAllocationLogic {
	return &SaveOperatorGameProviderAllocationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveOperatorGameProviderAllocationLogic) SaveOperatorGameProviderAllocation(req *types.SaveOperatorGameProviderAllocationRequest) (resp *types.SaveOperatorGameProviderAllocationResponse, err error) {
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameProviderServiceClient().SaveOperatorGameProviderAllocation(l.ctx, &platform_game.SaveOperatorGameProviderAllocationRequest{
		OpCode: req.OpCode,
		Items:  convertToRpcSaveOperatorGameProviderAllocationInfo(req.Items),
	})
	if err != nil {
		return nil, err
	}
	return &types.SaveOperatorGameProviderAllocationResponse{
		Total:   rpcResp.Total,
		Created: rpcResp.Created,
		Deleted: rpcResp.Deleted,
		Exist:   rpcResp.Exist,
		Failed:  rpcResp.Failed,
	}, nil
}

func convertToRpcSaveOperatorGameProviderAllocationInfo(items []types.SaveOperatorGameProviderAllocationItem) []*platform_game.SaveOperatorGameProviderAllocationInfo {
	rpcItems := make([]*platform_game.SaveOperatorGameProviderAllocationInfo, len(items))
	for i, item := range items {
		rpcItems[i] = &platform_game.SaveOperatorGameProviderAllocationInfo{
			Code:        item.Code,
			CheckStatus: item.CheckStatus,
		}
	}
	return rpcItems
}
