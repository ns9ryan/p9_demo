package gameservicelogic

import (
	"context"
	"fmt"

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

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetGame] DAO Manager not available")
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	gameRecord, err := l.svcCtx.DAOManager.Game.GetGameByID(l.ctx, in.GetId())
	if err != nil {
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

	if svcCtx == nil || svcCtx.DAOManager == nil {
		return nil, fmt.Errorf("DAO Manager not available")
	}

	// 查询分类信息
	categoryRecord, err := svcCtx.DAOManager.GameCategory.GetGameCategoryBySourceId(ctx, gameRecord.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("query category failed: categoryID=%d: %w", gameRecord.CategoryID, err)
	}
	ext.CategoryCode = categoryRecord.SourceCategoryCode

	// 查询供应商信息
	providerRecord, err := svcCtx.DAOManager.GameProvider.GetGameProviderBySourceId(ctx, gameRecord.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("query provider failed: providerID=%d: %w", gameRecord.ProviderID, err)
	}
	ext.ProviderCode = providerRecord.SourceProviderCode

	// 查询渠道信息（如果存在）
	if gameRecord.ChannelID != 0 {
		channelRecord, err := svcCtx.DAOManager.GameChannel.GetGameChannelBySourceId(ctx, gameRecord.ChannelID)
		if err != nil {
			return nil, fmt.Errorf("query channel failed: channelID=%d: %w", gameRecord.ChannelID, err)
		}
		ext.ChannelCode = channelRecord.SourceChannelCode
	}

	// 查询游戏货币信息 - 需要通过GameCurrencyDAO
	gameCurrencyRecords, err := svcCtx.DAOManager.GameCurrency.GetAllGameCurrency(ctx, gameRecord.SourceID, 0)
	if err != nil {
		return nil, fmt.Errorf("query game currency failed: %w", err)
	}

	gameCurrencyInfo := []GameCurrencyInfo{}
	for _, gameCurrencyRecord := range gameCurrencyRecords {
		currencyRecord, err := svcCtx.DAOManager.Currency.GetCurrencyByID(ctx, gameCurrencyRecord.CurrencyID)
		if err != nil {
			fmt.Errorf("query currency failed: currencyID=%d: %w", gameCurrencyRecord.CurrencyID, err)
			continue
		}
		// 需要查询Currency表获取name_key
		gameCurrencyInfo = append(gameCurrencyInfo, GameCurrencyInfo{
			CurrencyID:      gameCurrencyRecord.CurrencyID,
			CurrencyNameKey: currencyRecord.NameKey,
		})
	}

	ext.GameCurrencyInfo = utils.JSON(gameCurrencyInfo)
	return ext, nil
}

type GameCurrencyInfo struct {
	CurrencyID      int64  `json:"currency_id" comment:"币种ID"`
	CurrencyNameKey string `json:"currency_name_key" comment:"币种名称Key"`
}
