package gameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameLogic {
	return &GetGameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取单个游戏详情
func (l *GetGameLogic) GetGame(in *platformgame.GetGameRequest) (*platformgame.GetGameResp, error) {
	l.Infof("[RPC GetGame] received request: id=%d", in.GetId())

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGame] Database not available")
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	gameRecord := &ent.Game{}
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.GetId()).
		First(gameRecord).Error; err != nil {
		l.Errorf("[RPC GetGame] query failed: %v", err)
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get game: " + err.Error(),
		}, nil
	}
	// 查询关联的分类、供应商和渠道信息
	ext, err := GetGameExtInfo(l.ctx, l.svcCtx, gameRecord)
	if err != nil {
		l.Errorf("[RPC GetGame] query extended game info failed: %v", err)
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get extended game info: " + err.Error(),
		}, nil
	}
	return &platformgame.GetGameResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.GameModelToProto(gameRecord, ext),
	}, nil
}

func GetGameExtInfo(ctx context.Context, svcCtx *svc.ServiceContext, gameRecord *ent.Game) (*logic.GameInfoExt, error) {
	ext := &logic.GameInfoExt{}
	categoryRecord := &ent.GameCategory{}
	if err := svcCtx.DB.WithContext(ctx).
		Where("source_id = ? AND deleted_at IS NULL", gameRecord.CategoryId).
		First(categoryRecord).Error; err != nil {
		return nil, err
	}
	ext.CategoryNameI18N = categoryRecord.NameI18n

	providerRecord := &ent.GameProvider{}
	if err := svcCtx.DB.WithContext(ctx).
		Where("source_id = ? AND deleted_at IS NULL", gameRecord.ProviderId).
		First(providerRecord).Error; err != nil {
		return nil, err
	}
	ext.ProviderNameI18N = providerRecord.NameI18n

	channelRecord := &ent.GameChannel{}
	if gameRecord.ChannelId.Int64 != 0 {
		if err := svcCtx.DB.WithContext(ctx).
			Where("source_id = ? AND deleted_at IS NULL", gameRecord.ChannelId.Int64).
			First(channelRecord).Error; err != nil {
			return nil, err
		}
	}
	ext.ChannelNameI18N = channelRecord.NameI18n

	gameCurrencyRecords := []*ent.GameCurrency{}
	if err := svcCtx.DB.WithContext(ctx).
		Where("game_id = ? AND deleted_at IS NULL", gameRecord.SourceId).
		Find(&gameCurrencyRecords).Error; err != nil {
		return nil, err
	}
	gameCurrencyInfo := []GameCurrencyInfo{}
	for _, currency := range gameCurrencyRecords {
		if currency == nil {
			continue
		}
		sysCurrencyRecords := &ent.SysCurrency{}
		if err := svcCtx.DB.WithContext(ctx).
			Where("id = ?", currency.CurrencyId).
			First(sysCurrencyRecords).Error; err != nil {
			return nil, err
		}
		gameCurrencyInfo = append(gameCurrencyInfo, GameCurrencyInfo{
			CurrencyID:       currency.CurrencyId,
			CurrencyNameI18n: sysCurrencyRecords.NameI18n,
		})
	}
	ext.GameCurrencyInfo = utils.JSON(gameCurrencyInfo)
	return ext, nil
}

type GameCurrencyInfo struct {
	CurrencyID       int64  `json:"currency_id" comment:"币种ID"`
	CurrencyNameI18n string `json:"currency_name_i18n" comment:"币种名称（多语言JSON）"`
}
