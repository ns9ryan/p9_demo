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

type BatchUpdateOperatorGameChannelStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量修改分站游戏渠道状态
func NewBatchUpdateOperatorGameChannelStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameChannelStatusLogic {
	return &BatchUpdateOperatorGameChannelStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateOperatorGameChannelStatusLogic) BatchUpdateOperatorGameChannelStatus(req *types.BatchUpdateOperatorGameChannelStatusRequest) (resp *types.BatchUpdateOperatorGameChannelStatusResponse, err error) {
	// 调用 RPC 批量修改分站游戏渠道状态
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameChannelServiceClient().BatchUpdateOperatorGameChannelStatus(l.ctx, &platform_game.BatchUpdateOperatorGameChannelStatusRequest{
		Ids:    req.Ids,
		Status: req.Status,
	})
	if err != nil {
		l.Logger.Error("BatchUpdateOperatorGameChannelStatus error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchUpdateOperatorGameChannelStatusResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
