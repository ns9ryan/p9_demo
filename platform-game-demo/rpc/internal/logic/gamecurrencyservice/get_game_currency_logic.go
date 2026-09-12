package gamecurrencyservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
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
func (l *GetGameCurrencyLogic) GetGameCurrency(in *platformgame.GetGameCurrencyRequest) (*platformgame.GetGameCurrencyResp, error) {
	l.Infof("[RPC GetGameCurrency] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameCurrency] Database not available")
		return &platformgame.GetGameCurrencyResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	currency := &ent.GameCurrency{}
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(currency).Error; err != nil {
		l.Errorf("[RPC GetGameCurrency] query failed: %v", err)
		return &platformgame.GetGameCurrencyResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get currency: " + err.Error(),
		}, nil
	}

	sysCurrencyMap := GetSysCurrencyMap(l.ctx, l.svcCtx)

	gameRecord, err := GetGameRecord(l.ctx, l.svcCtx, currency.GameId)
	if err != nil {
		l.Errorf("[RPC GetGameCurrency] query game failed: %v", err)
		return &platformgame.GetGameCurrencyResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get game: " + err.Error(),
		}, nil
	}

	l.Infof("[RPC GetGameCurrency] query result: id=%d, game_id=%d, currency_id=%d, status=%d, deleted_at=%v",
		currency.Id, currency.GameId, currency.CurrencyId, currency.Status, currency.DeletedAt)

	return &platformgame.GetGameCurrencyResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.CurrencyModelToProto(currency, gameRecord, sysCurrencyMap),
	}, nil
}

func GetSysCurrencyMap(ctx context.Context, svcCtx *svc.ServiceContext) map[int64]interface{} {
	sysCurrencyRecords := []*ent.SysCurrency{}
	err := svcCtx.DB.WithContext(ctx).Model(&ent.SysCurrency{}).
		Select("id", "currency_code", "name_i18n").
		Find(&sysCurrencyRecords).Error
	if err != nil {
		// handle error appropriately, e.g., log and return an empty map
		logx.Errorf("[RPC GetGameCurrencyList] query system currency failed: %v", err)
		return map[int64]interface{}{}
	}
	currencyMap := make(map[int64]interface{})
	for _, sc := range sysCurrencyRecords {
		currencyMap[sc.Id] = map[string]string{
			"currency_code": sc.CurrencyCode,
			"name_i18n":     sc.NameI18n,
		}
	}
	return currencyMap
}

func GetGameRecord(ctx context.Context, svcCtx *svc.ServiceContext, gameId int64) (*ent.Game, error) {
	gameRecord := &ent.Game{}
	if err := svcCtx.DB.WithContext(ctx).
		Select("game_code", "name_i18n").
		Where("source_id = ? AND deleted_at IS NULL", gameId).
		First(gameRecord).Error; err != nil {
		return &ent.Game{}, err
	}
	return gameRecord, nil
}
