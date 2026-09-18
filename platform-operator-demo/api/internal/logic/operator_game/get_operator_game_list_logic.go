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

type GetOperatorGameListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分站游戏列表
func NewGetOperatorGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameListLogic {
	return &GetOperatorGameListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOperatorGameListLogic) GetOperatorGameList(req *types.GetOperatorGameListRequest) (resp *types.GetOperatorGameListResponse, err error) {
	// 调用 RPC 获取分站游戏列表
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameServiceClient().GetOperatorGameList(l.ctx, &platform_game.GetOperatorGameListRequest{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		OpCode:   req.OpCode,
		GameCode: req.GameCode,
		Status:   req.Status,
	})
	if err != nil {
		l.Logger.Error("GetOperatorGameList error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.GetOperatorGameListResponse{
		Total:    rpcResp.Total,
		Page:     rpcResp.Page,
		PageSize: rpcResp.PageSize,
	}

	// 转换数据
	for _, item := range rpcResp.Items {
		resp.Items = append(resp.Items, types.OperatorGameInfo{
			Id:        item.Id,
			OpCode:    item.OpCode,
			GameCode:  item.GameCode,
			Name:      item.Name,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return resp, nil
}
