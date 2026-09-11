// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package provider

import (
	"context"

	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameProviderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameProviderListLogic {
	return &GameProviderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameProviderListLogic) GameProviderList(req *types.GameProviderListReq) (resp *types.GameProviderListResp, err error) {
	l.Infof("[API GameProviderList] received req: page=%d, page_size=%d", req.Page, req.PageSize)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameProviderList] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.GetGameProviderListRequest{
		Page:         int32(req.Page),
		PageSize:     int32(req.PageSize),
		ProviderCode: req.ProviderCode,
		Name:         req.Name,
		Status:       int32(req.Status),
		IsDeleted:    int32(req.IsDeleted),
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
	}

	client := l.svcCtx.GrpcClient.GetGameProviderServiceClient()
	grpcResp, err := client.GetGameProviderList(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameProviderList] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil {
		l.Error("[API GameProviderList] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameProviderList] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	items := make([]types.GameProviderResp, 0, len(grpcResp.Data))
	for _, item := range grpcResp.Data {
		items = append(items, *logic.ProviderProtoToResponse(item))
	}

	resp = &types.GameProviderListResp{
		List:  items,
		Total: grpcResp.Total,
	}

	l.Infof("[API GameProviderList] success: total=%d", grpcResp.Total)
	return resp, nil
}
