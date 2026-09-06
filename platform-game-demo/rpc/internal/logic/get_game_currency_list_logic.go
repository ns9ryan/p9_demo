package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/pkg/game"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
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
	l.Infof("GetGameCurrencyList called with page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx.DB == nil {
		l.Errorf("Database not available")
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	// 构建查询参数
	page := int64(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int64(in.PageSize)
	if pageSize <= 0 {
		pageSize = 15
	}

	params := &game.GameCurrencyListParams{
		Page:       page,
		PageSize:   pageSize,
		GameID:     in.GetGameId(),
		CurrencyID: in.GetCurrencyId(),
		Status:     int16(in.GetStatus()),
		IsDeleted:  int8(in.GetIsDeleted()),
		SortBy:     in.GetSortBy(),
		SortOrder:  in.GetSortOrder(),
	}

	// 调用数据库查询函数
	result, err := game.GameCurrencyListDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("GameCurrencyListDB failed: %v", err)
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to fetch currencies: " + err.Error(),
		}, nil
	}

	// 转换结果为 protobuf 消息
	currencies := make([]*platformgame.CurrencyInfo, 0, len(result.Currencies))
	for _, currencyWithMeta := range result.Currencies {
		if currencyWithMeta == nil {
			continue
		}

		// 查询 Game 表获取 game_name_i18n
		var game model.Game
		gameNameI18n := ""
		if currencyWithMeta.GameID > 0 {
			if err := l.svcCtx.DB.WithContext(l.ctx).
				Where("id = ? AND deleted_at IS NULL", currencyWithMeta.GameID).
				First(&game).Error; err == nil {
				gameNameI18n = string(utils.MustMarshalJSON(game.NameI18n))
			}
		}

		// 查询 sys_currency 表获取 currency_name_i18n
		var currencyData struct {
			NameI18n string
		}
		currencyNameI18n := ""
		if currencyWithMeta.CurrencyID > 0 {
			if err := l.svcCtx.DB.WithContext(l.ctx).
				Table("sys_currency").
				Where("id = ?", currencyWithMeta.CurrencyID).
				Select("name_i18n").
				Scan(&currencyData).Error; err == nil && currencyData.NameI18n != "" {
				currencyNameI18n = currencyData.NameI18n
			}
		}

		// 使用 CurrencyCode 作为名称（GameCurrencyWithMeta 中只包含 code，没包含 NameI18n�?
		currencies = append(currencies, &platformgame.CurrencyInfo{
			Id:               currencyWithMeta.CurrencyID,
			Code:             currencyWithMeta.CurrencyCode,
			Name:             currencyWithMeta.CurrencyCode,
			GameId:           currencyWithMeta.GameID,
			CurrencyId:       currencyWithMeta.CurrencyID,
			GameNameI18N:     gameNameI18n,
			CurrencyNameI18N: currencyNameI18n,
		})
	}

	return &platformgame.GetGameCurrencyListResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data:    currencies,
		Total:   result.Total,
	}, nil
}
