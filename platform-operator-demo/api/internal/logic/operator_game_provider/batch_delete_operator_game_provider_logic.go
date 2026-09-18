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

type BatchDeleteOperatorGameProviderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除分站游戏提供商
func NewBatchDeleteOperatorGameProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperatorGameProviderLogic {
	return &BatchDeleteOperatorGameProviderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteOperatorGameProviderLogic) BatchDeleteOperatorGameProvider(req *types.BatchDeleteOperatorGameProviderRequest) (resp *types.BatchDeleteOperatorGameProviderResponse, err error) {
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameProviderServiceClient().BatchDeleteOperatorGameProvider(l.ctx, &platform_game.BatchDeleteOperatorGameProviderRequest{
		Ids: req.Ids,
	})
	if err != nil {
		l.Logger.Error("BatchDeleteOperatorGameProvider error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchDeleteOperatorGameProviderResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
