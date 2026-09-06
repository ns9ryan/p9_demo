// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameProviderGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameProviderGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameProviderGetLogic {
	return &GameProviderGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameProviderGetLogic) GameProviderGet(req *types.GameProviderGetReq) (resp *types.GameProviderResp, err error) {
	l.Infof("[API GameProviderGet] received req: id=%d", req.ID)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameProviderGet] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.GetGameProviderRequest{
		Id: req.ID,
	}

	client := l.svcCtx.GrpcClient.GetPlatformGameServiceClient()
	grpcResp, err := client.GetGameProvider(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameProviderGet] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil || len(grpcResp.Data) == 0 {
		l.Error("[API GameProviderGet] gRPC response is nil or empty")
		return nil, fmt.Errorf("gRPC response is nil or empty")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameProviderGet] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	item := grpcResp.Data[0]
	resp = &types.GameProviderResp{
		ID:           item.Id,
		ProviderCode: item.Code,
		NameI18n:     item.Name,
		Status:       int16(item.Status),
	}

	l.Infof("[API GameProviderGet] success: id=%d", item.Id)
	return resp, nil
}
