// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGameAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取游戏资源分配列表
func NewListGameAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGameAllocationsLogic {
	return &ListGameAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGameAllocationsLogic) ListGameAllocations(req *types.ListGameAllocationsRequest) (resp *types.ListGameAllocationsResponse, err error) {
	// 调用 RPC 获取游戏分配列表
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameAllocationServiceClient().GetOperatorGameAllocationList(l.ctx, &platform_game.GetOperatorGameAllocationListRequest{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		OpCode:   req.OpCode,
		Status:   req.Status,
	})
	if err != nil {
		l.Logger.Error("GetOperatorGameAllocationList error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.ListGameAllocationsResponse{
		Code:     rpcResp.Code,
		Message:  rpcResp.Message,
		Total:    rpcResp.Total,
		Page:     rpcResp.Page,
		PageSize: rpcResp.PageSize,
	}

	// 转换数据
	for _, item := range rpcResp.Items {
		resp.Items = append(resp.Items, types.OperatorGameAllocationInfo{
			Id:                item.Id,
			OpCode:            item.OpCode,
			OpName:            item.OpName,
			GameCount:         item.GameCount,
			GameCategoryCount: item.GameCategoryCount,
			GameProviderCount: item.GameProviderCount,
			GameChannelCount:  item.GameChannelCount,
			UpdatedAt:         item.UpdatedAt,
		})
	}

	return resp, nil
}
