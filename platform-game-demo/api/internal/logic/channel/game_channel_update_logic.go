// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package channel

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/constant"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type GameChannelUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameChannelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameChannelUpdateLogic {
	return &GameChannelUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameChannelUpdateLogic) GameChannelUpdate(req *types.GameChannelUpdateReq) (resp *types.GameChannelResp, err error) {
	if l.svcCtx == nil || l.svcCtx.GameGrpcClient == nil {
		logger.Error("[API GameChannelUpdate] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platform_game.UpdateGameChannelRequest{
		Id:          req.ID,
		SortNo:      req.SortNo,
		Status:      int32(req.Status),
		ForceLogout: req.ForceLogout,
	}

	grpcResp, err := l.svcCtx.GameGrpcClient.GetGameChannelServiceClient().UpdateGameChannel(l.ctx, grpcReq)
	if err != nil {
		logger.Errorf("[API GameChannelUpdate] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp == nil {
		logger.Error("[API GameChannelUpdate] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		logger.Errorf("[API GameChannelUpdate] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	resp = logic.ChannelProtoToResponse(l.ctx, grpcResp.Data)

	logger.Infof("[API GameChannelUpdate] success: id=%d", req.ID)
	return resp, nil
}
