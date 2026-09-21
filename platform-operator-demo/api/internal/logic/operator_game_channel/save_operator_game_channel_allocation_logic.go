// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_channel

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveOperatorGameChannelAllocationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存分站游戏渠道分配
func NewSaveOperatorGameChannelAllocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveOperatorGameChannelAllocationLogic {
	return &SaveOperatorGameChannelAllocationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveOperatorGameChannelAllocationLogic) SaveOperatorGameChannelAllocation(req *types.SaveOperatorGameChannelAllocationRequest) (resp *types.SaveOperatorGameChannelAllocationResponse, err error) {
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameChannelServiceClient().SaveOperatorGameChannelAllocation(l.ctx, &platform_game.SaveOperatorGameChannelAllocationRequest{
		OpCode: req.OpCode,
		Items:  convertToRpcSaveOperatorGameChannelAllocationInfo(req.Items),
	})
	if err != nil {
		return nil, err
	}
	return &types.SaveOperatorGameChannelAllocationResponse{
		Total:   rpcResp.Total,
		Created: rpcResp.Created,
		Deleted: rpcResp.Deleted,
		Exist:   rpcResp.Exist,
		Failed:  rpcResp.Failed,
	}, nil
}

func convertToRpcSaveOperatorGameChannelAllocationInfo(items []types.SaveOperatorGameChannelAllocationItem) []*platform_game.SaveOperatorGameChannelAllocationInfo {
	rpcItems := make([]*platform_game.SaveOperatorGameChannelAllocationInfo, len(items))
	for i, item := range items {
		rpcItems[i] = &platform_game.SaveOperatorGameChannelAllocationInfo{
			Code:        item.Code,
			CheckStatus: item.CheckStatus,
		}
	}
	return rpcItems
}
