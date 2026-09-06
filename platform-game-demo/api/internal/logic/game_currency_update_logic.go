// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/logger"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type GameCurrencyUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCurrencyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCurrencyUpdateLogic {
	return &GameCurrencyUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameCurrencyUpdateLogic) GameCurrencyUpdate(req *types.GameCurrencyUpdateReq) (resp *types.GameCurrencyResp, err error) {
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		logger.Error("[API GameCurrencyUpdate] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.UpdateGameCurrencyRequest{
		Id:          req.ID,
		Status:      int32(req.Status),
		ForceLogout: req.ForceLogout,
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetPlatformGameServiceClient().UpdateGameCurrency(l.ctx, grpcReq)
	if err != nil {
		logger.Errorf("[API GameCurrencyUpdate] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp == nil {
		logger.Error("[API GameCurrencyUpdate] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		logger.Errorf("[API GameCurrencyUpdate] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	resp = &types.GameCurrencyResp{
		ID:     req.ID,
		Status: req.Status,
	}

	logger.Infof("[API GameCurrencyUpdate] success: id=%d", req.ID)
	return resp, nil
}
