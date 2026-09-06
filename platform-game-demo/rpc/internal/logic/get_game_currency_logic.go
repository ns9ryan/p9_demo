package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameCurrencyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameCurrencyLogic {
	return &GetGameCurrencyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取单个游戏货币
func (l *GetGameCurrencyLogic) GetGameCurrency(in *platformgame.GetGameCurrencyRequest) (*platformgame.GetGameCurrencyListResp, error) {
	l.Infof("[RPC GetGameCurrency] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameCurrency] Database not available")
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
			Data:    nil,
		}, nil
	}

	var currency model.GameCurrency
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&currency).Error; err != nil {
		l.Errorf("[RPC GetGameCurrency] query failed: %v", err)
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get currency: " + err.Error(),
			Data:    nil,
		}, nil
	}

	l.Infof("[RPC GetGameCurrency] query result: id=%d, game_id=%d, currency_id=%d, status=%d, deleted_at=%v",
		currency.ID, currency.GameID, currency.CurrencyID, currency.Status, currency.DeletedAt)

	var game model.Game
	gameNameI18n := ""
	if currency.GameID > 0 {
		if err := l.svcCtx.DB.WithContext(l.ctx).
			Where("id = ? AND deleted_at IS NULL", currency.GameID).
			First(&game).Error; err == nil {
			gameNameI18n = string(utils.MustMarshalJSON(game.NameI18n))
		}
	}

	// 查询 sys_currency 表获取 currency_name_i18n
	var currencyData struct {
		NameI18n string
	}
	currencyNameI18n := ""
	if currency.CurrencyID > 0 {
		if err := l.svcCtx.DB.WithContext(l.ctx).
			Table("sys_currency").
			Where("id = ?", currency.CurrencyID).
			Select("name_i18n").
			Scan(&currencyData).Error; err == nil && currencyData.NameI18n != "" {
			currencyNameI18n = currencyData.NameI18n
		}
	}

	item := &platformgame.CurrencyInfo{
		Id:               currency.ID,
		GameId:           currency.GameID,
		CurrencyId:       currency.CurrencyID,
		GameNameI18N:     gameNameI18n,
		CurrencyNameI18N: currencyNameI18n,
	}

	return &platformgame.GetGameCurrencyListResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    []*platformgame.CurrencyInfo{item},
	}, nil
}
