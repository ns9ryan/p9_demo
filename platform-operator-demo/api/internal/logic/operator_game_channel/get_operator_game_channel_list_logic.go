// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_game_channel

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分站游戏渠道列表
func NewGetOperatorGameChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameChannelListLogic {
	return &GetOperatorGameChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOperatorGameChannelListLogic) GetOperatorGameChannelList(req *types.GetOperatorGameChannelListRequest) (resp *types.GetOperatorGameChannelListResponse, err error) {
	// 调用 RPC 获取分站游戏渠道列表
	rpcResp, err := l.svcCtx.GameGrpcClient.GetOperatorGameChannelServiceClient().GetOperatorGameChannelList(l.ctx, &platform_game.GetOperatorGameChannelListRequest{
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
		OpCode:      req.OpCode,
		ChannelCode: req.ChannelCode,
		Status:      req.Status,
	})
	if err != nil {
		l.Logger.Error("GetOperatorGameChannelList error:", err)
		return nil, err
	}

	// 构建响应
	resp = &types.GetOperatorGameChannelListResponse{
		Total:    rpcResp.Total,
		Page:     rpcResp.Page,
		PageSize: rpcResp.PageSize,
	}

	// 转换数据，並使用 corei18n.TG 获取 name
	for _, item := range rpcResp.Items {
		name := i18n.TG(l.ctx, i18n.CodePlatform, "game", fmt.Sprintf("game.channel.%s.name", item.ChannelCode))
		resp.Items = append(resp.Items, types.OperatorGameChannelInfo{
			Id:          item.Id,
			OpCode:      item.OpCode,
			ChannelCode: item.ChannelCode,
			Name:        name,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return resp, nil
}
