// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package currency

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

type GameCurrencyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCurrencyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCurrencyUpdateLogic {
	return &GameCurrencyUpdateLogic{
		Logger: logx.WithContext(ctx),
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

	grpcResp, err := l.svcCtx.GrpcClient.GetGameCurrencyServiceClient().UpdateGameCurrency(l.ctx, grpcReq)
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

	resp = logic.CurrencyProtoToResponse(grpcResp.Data)

	logger.Infof("[API GameCurrencyUpdate] success: id=%d", req.ID)
	return resp, nil
}
