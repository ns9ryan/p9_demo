package gamecurrencyservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type GetGameCurrencyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameCurrencyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameCurrencyListLogic {
	return &GetGameCurrencyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏货币列表
func (l *GetGameCurrencyListLogic) GetGameCurrencyList(in *platformgame.GetGameCurrencyListRequest) (*platformgame.GetGameCurrencyListResp, error) {
	l.Infof("[RPC GetGameCurrencyList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetGameCurrencyList] DAO Manager not available")
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	l.Infof("[RPC GetGameCurrencyList] query params: isDeleted=%d, status=%d, gameId=%d, currencyId=%d, offset=%d, pageSize=%d",
		in.GetIsDeleted(), in.GetStatus(), in.GetGameId(), in.GetCurrencyId(), offset, pageSize)

	currencies, total, err := l.svcCtx.DAOManager.GameCurrency.GetGameCurrencyList(
		l.ctx,
		in.GetIsDeleted(),
		in.GetStatus(),
		in.GetGameId(),
		in.GetCurrencyId(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetGameCurrencyList] DAO query failed: %v", err)
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	l.Infof("[RPC GetGameCurrencyList] DAO returned: data_count=%d, total=%d", len(currencies), total)

	sysCurrencyMap, err := l.svcCtx.DAOManager.Currency.GetCurrencyMap(l.ctx)
	if err != nil {
		l.Errorf("[RPC GetGameCurrencyList] query system currency failed: %v", err)
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get system currency: " + err.Error(),
		}, nil
	}
	currenciesProto := make([]*platformgame.GameCurrencyInfo, 0, len(currencies))
	for _, currency := range currencies {
		if currency == nil {
			continue
		}
		gameRecord, err := GetGameRecord(l.ctx, l.svcCtx, currency.GameID)
		if err != nil {
			l.Errorf("[RPC GetGameCurrencyList] query game failed: %v", err)
		}
		currenciesProto = append(currenciesProto, logic.CurrencyModelToProto(currency, gameRecord, sysCurrencyMap))
	}

	l.Infof("[RPC GetGameCurrencyList] success: total=%d", total)
	return &platformgame.GetGameCurrencyListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Data:     currenciesProto,
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
