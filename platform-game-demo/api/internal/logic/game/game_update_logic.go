// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package game

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

type GameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameUpdateLogic {
	return &GameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameUpdateLogic) GameUpdate(req *types.GameUpdateReq) (resp *types.GameResp, err error) {
	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		logger.Error("[API GameUpdate] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	var supportsEmbed int32 = 0
	var supportsRedirect int32 = 0
	if req.SupportsEmbed != nil {
		if *req.SupportsEmbed == true {
			supportsEmbed = 1
		}
		if *req.SupportsEmbed == false {
			supportsEmbed = 2
		}
	}
	if req.SupportsRedirect != nil {
		if *req.SupportsRedirect == true {
			supportsRedirect = 1
		}
		if *req.SupportsRedirect == false {
			supportsRedirect = 2
		}
	}

	grpcReq := &platformgame.UpdateGameRequest{
		Id:               req.ID,
		NameI18N:         req.NameI18n,
		SortNo:           req.SortNo,
		Status:           int32(req.Status),
		ImageUrl:         req.ImageUrl,
		SupportsEmbed:    supportsEmbed,
		SupportsRedirect: supportsRedirect,
		ForceLogout:      req.ForceLogout,
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetGameServiceClient().UpdateGame(l.ctx, grpcReq)
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

	resp = logic.GameProtoToResponse(grpcResp.Data)

	logger.Infof("[API GameUpdate] success: id=%d", req.ID)
	return resp, nil
}
