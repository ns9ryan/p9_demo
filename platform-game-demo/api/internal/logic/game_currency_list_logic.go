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

type GameCurrencyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameCurrencyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCurrencyListLogic {
	return &GameCurrencyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameCurrencyListLogic) GameCurrencyList(req *types.GameCurrencyListReq) (resp *types.GameCurrencyListResp, err error) {
	l.Infof("[API GameCurrencyList] received req: page=%d, page_size=%d", req.Page, req.PageSize)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Error("[API GameCurrencyList] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.GetGameCurrencyListRequest{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),

		IsDeleted: int32(req.IsDeleted),
	}

	client := l.svcCtx.GrpcClient.GetPlatformGameServiceClient()
	grpcResp, err := client.GetGameCurrencyList(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameCurrencyList] gRPC call failed: %v", err)
		return nil, err
	}

	if grpcResp == nil {
		l.Error("[API GameCurrencyList] gRPC response is nil")
		return nil, fmt.Errorf("gRPC response is nil")
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameCurrencyList] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	items := make([]types.GameCurrencyResp, 0, len(grpcResp.Data))
	for _, item := range grpcResp.Data {
		items = append(items, types.GameCurrencyResp{
			ID:               item.Id,
			GameID:           item.GameId,
			CurrencyID:       item.CurrencyId,
			CurrencyCode:     item.Code,
			GameNameI18n:     item.GameNameI18N,
			CurrencyNameI18n: item.CurrencyNameI18N,
		})
	}

	resp = &types.GameCurrencyListResp{
		Items: items,
		Total: grpcResp.Total,
	}

	l.Infof("[API GameCurrencyList] success: total=%d", grpcResp.Total)
	return resp, nil
}
