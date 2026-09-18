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

type BatchCreateOperatorGameChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建分站游戏渠道
func NewBatchCreateOperatorGameChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameChannelLogic {
	return &BatchCreateOperatorGameChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCreateOperatorGameChannelLogic) BatchCreateOperatorGameChannel(req *types.BatchCreateOperatorGameChannelRequest) (resp *types.BatchCreateOperatorGameChannelResponse, err error) {
	// 构建 RPC 请求
	var items []*platform_game.BatchCreateOperatorGameChannelItemRequest
	for _, item := range req.Items {
		items = append(items, &platform_game.BatchCreateOperatorGameChannelItemRequest{
			OpCode:      item.OpCode,
			ChannelCode: item.ChannelCode,
			Status:      item.Status,
		})
	}

	// 调用 RPC 批量创建分站游戏渠道
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameChannelServiceClient().BatchCreateOperatorGameChannel(l.ctx, &platform_game.BatchCreateOperatorGameChannelRequest{
		Items: items,
	})
	if err != nil {
		l.Logger.Error("BatchCreateOperatorGameChannel error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchCreateOperatorGameChannelResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
