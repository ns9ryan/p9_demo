// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_provider

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameProviderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分站游戏提供商列表
func NewGetOperatorGameProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameProviderListLogic {
	return &GetOperatorGameProviderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOperatorGameProviderListLogic) GetOperatorGameProviderList(req *types.GetOperatorGameProviderListRequest) (resp *types.GetOperatorGameProviderListResponse, err error) {
	// 调用 RPC 获取分站游戏提供商列表
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameProviderServiceClient().GetOperatorGameProviderList(l.ctx, &platform_game.GetOperatorGameProviderListRequest{
		Page:         int32(req.Page),
		PageSize:     int32(req.PageSize),
		OpCode:       req.OpCode,
		ProviderCode: req.ProviderCode,
		Status:       req.Status,
	})
	if err != nil {
		l.Logger.Error("GetOperatorGameProviderList error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.GetOperatorGameProviderListResponse{
		Total:    rpcResp.Total,
		Page:     rpcResp.Page,
		PageSize: rpcResp.PageSize,
	}

	// 转换数据，並使用 corei18n.TG 获取 name
	for _, item := range rpcResp.Items {
		name := i18n.TG(l.ctx, i18n.CodePlatform, "game", fmt.Sprintf("game.provider.%s.name", item.ProviderCode))
		resp.Items = append(resp.Items, types.OperatorGameProviderInfo{
			Id:           item.Id,
			OpCode:       item.OpCode,
			ProviderCode: item.ProviderCode,
			Name:         name,
			Status:       item.Status,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	return resp, nil
}
