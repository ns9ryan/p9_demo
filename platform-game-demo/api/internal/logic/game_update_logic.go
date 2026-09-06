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

type GameUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameUpdateLogic {
	return &GameUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameUpdateLogic) GameUpdate(req *types.GameUpdateReq) (resp *types.GameResp, err error) {
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		logger.Error("[API GameUpdate] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.UpdateGameRequest{
		Id:          req.ID,
		NameI18N:    req.NameI18n,
		SortNo:      req.SortNo,
		Status:      int32(req.Status),
		ForceLogout: req.ForceLogout,
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetPlatformGameServiceClient().UpdateGame(l.ctx, grpcReq)
	if err != nil {
		logger.Errorf("[API GameUpdate] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp == nil {
		logger.Error("[API GameUpdate] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		logger.Errorf("[API GameUpdate] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	resp = &types.GameResp{
		ID:       req.ID,
		GameCode: grpcResp.Data.Code,
		NameI18n: grpcResp.Data.Name,
		Status:   int16(grpcResp.Data.Status),
	}

	logger.Infof("[API GameUpdate] success: id=%d", req.ID)
	return resp, nil
}
