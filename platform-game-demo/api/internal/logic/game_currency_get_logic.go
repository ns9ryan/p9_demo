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

type GameCurrencyGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCurrencyGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCurrencyGetLogic {
	return &GameCurrencyGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameCurrencyGetLogic) GameCurrencyGet(req *types.GameCurrencyGetReq) (resp *types.GameCurrencyResp, err error) {
	l.Infof("[API GameCurrencyGet] received req: id=%d", req.ID)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameCurrencyGet] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.GetGameCurrencyRequest{
		Id: req.ID,
	}

	client := l.svcCtx.GrpcClient.GetPlatformGameServiceClient()
	grpcResp, err := client.GetGameCurrency(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameCurrencyGet] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil || len(grpcResp.Data) == 0 {
		l.Error("[API GameCurrencyGet] gRPC response is nil or empty")
		return nil, fmt.Errorf("gRPC response is nil or empty")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameCurrencyGet] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	item := grpcResp.Data[0]
	l.Infof("[API GameCurrencyGet] gRPC response: id=%d, game_id=%d, currency_id=%d, game_name_i18n=%s", item.Id, item.GameId, item.CurrencyId, item.GameNameI18N)

	resp = &types.GameCurrencyResp{
		ID:               item.Id,
		GameID:           item.GameId,
		CurrencyID:       item.CurrencyId,
		GameNameI18n:     item.GameNameI18N,
		CurrencyNameI18n: item.CurrencyNameI18N,
	}

	l.Infof("[API GameCurrencyGet] success: id=%d, game_id=%d, currency_id=%d", item.Id, item.GameId, item.CurrencyId)
	return resp, nil
}
