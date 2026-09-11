// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package provider

import (
	"context"

	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"fmt"

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

	client := l.svcCtx.GrpcClient.GetGameProviderServiceClient()
	grpcResp, err := client.GetGameProvider(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameProviderGet] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil || grpcResp.Data == nil {
		l.Error("[API GameProviderGet] gRPC response is nil or empty")
		return nil, fmt.Errorf("gRPC response is nil or empty")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameProviderGet] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	item := grpcResp.Data
	resp = logic.ProviderProtoToResponse(item)

	l.Infof("[API GameProviderGet] success: id=%d", item.Id)
	return resp, nil
}
