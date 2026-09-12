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

	"github.com/bytedance/gopkg/util/logger"
	"github.com/zeromicro/go-zero/core/logx"
)

type GameProviderUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameProviderUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameProviderUpdateLogic {
	return &GameProviderUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameProviderUpdateLogic) GameProviderUpdate(req *types.GameProviderUpdateReq) (resp *types.GameProviderResp, err error) {
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		logger.Error("[API GameProviderUpdate] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.UpdateGameProviderRequest{
		Id:          req.ID,
		NameI18N:    req.NameI18n,
		SortNo:      req.SortNo,
		Status:      int32(req.Status),
		LogoUrl:     req.LogoUrl,
		ForceLogout: req.ForceLogout,
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetGameProviderServiceClient().UpdateGameProvider(l.ctx, grpcReq)
	if err != nil {
		logger.Errorf("[API GameProviderUpdate] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp == nil {
		logger.Error("[API GameProviderUpdate] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		logger.Errorf("[API GameProviderUpdate] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	if grpcResp.Data == nil {
		logger.Error("[API GameProviderUpdate] gRPC response data is empty")
		return nil, fmt.Errorf("gRPC response data is empty")
	}

	resp = logic.ProviderProtoToResponse(grpcResp.Data)

	logger.Infof("[API GameProviderUpdate] success: id=%d", req.ID)
	return resp, nil
}
