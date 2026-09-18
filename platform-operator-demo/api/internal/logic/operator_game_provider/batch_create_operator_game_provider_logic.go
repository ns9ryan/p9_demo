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

type BatchCreateOperatorGameProviderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建分站游戏提供商
func NewBatchCreateOperatorGameProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameProviderLogic {
	return &BatchCreateOperatorGameProviderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCreateOperatorGameProviderLogic) BatchCreateOperatorGameProvider(req *types.BatchCreateOperatorGameProviderRequest) (resp *types.BatchCreateOperatorGameProviderResponse, err error) {
	// 构建 RPC 请求
	var items []*platform_game.BatchCreateOperatorGameProviderItemRequest
	for _, item := range req.Items {
		items = append(items, &platform_game.BatchCreateOperatorGameProviderItemRequest{
			OpCode:       item.OpCode,
			ProviderCode: item.ProviderCode,
			Status:       item.Status,
		})
	}

	// 调用 RPC 批量创建分站游戏提供商
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameProviderServiceClient().BatchCreateOperatorGameProvider(l.ctx, &platform_game.BatchCreateOperatorGameProviderRequest{
		Items: items,
	})
	if err != nil {
		l.Logger.Error("BatchCreateOperatorGameProvider error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchCreateOperatorGameProviderResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
