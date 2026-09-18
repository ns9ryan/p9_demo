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

type BatchCreateOperatorGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建分站游戏
func NewBatchCreateOperatorGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameLogic {
	return &BatchCreateOperatorGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCreateOperatorGameLogic) BatchCreateOperatorGame(req *types.BatchCreateOperatorGameRequest) (resp *types.BatchCreateOperatorGameResponse, err error) {
	// 构建 RPC 请求
	var items []*platform_game.BatchCreateOperatorGameItemRequest
	for _, item := range req.Items {
		items = append(items, &platform_game.BatchCreateOperatorGameItemRequest{
			OpCode:   item.OpCode,
			GameCode: item.GameCode,
			Name:     item.Name,
			Status:   item.Status,
		})
	}

	// 调用 RPC 批量创建分站游戏
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameServiceClient().BatchCreateOperatorGame(l.ctx, &platform_game.BatchCreateOperatorGameRequest{
		Items: items,
	})
	if err != nil {
		l.Logger.Error("BatchCreateOperatorGame error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.BatchCreateOperatorGameResponse{
		Total:   rpcResp.Total,
		Success: rpcResp.Success,
		Failed:  rpcResp.Failed,
	}

	return resp, nil
}
