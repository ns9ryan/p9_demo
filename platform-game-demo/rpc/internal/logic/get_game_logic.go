package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/pkg/game"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
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
	l.Infof("GetGame called with game_id=%d", in.GameId)

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Errorf("Database not available")
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	// 构建查询参数
	params := &game.GameGetParams{
		ID: in.GameId,
	}

	// 调用数据库查询函�?
	result, err := game.GameGetDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("GameGetDB failed: %v", err)
		return &platformgame.GetGameResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to fetch game: " + err.Error(),
		}, nil
	}

	if result == nil || result.Game == nil {
		return &platformgame.GetGameResp{
			Code:    constant.CodeNotFound,
			Message: "Game not found",
		}, nil
	}

	// 转换结果为 protobuf 消息，保留原始 name_i18n JSON
	g := result.Game
	nameI18nRaw := ""
	if g.NameI18n != nil {
		nameI18nRaw = string(utils.MustMarshalJSON(g.NameI18n))
	}
	sourceNameI18nRaw := ""
	if g.SourceNameI18n != nil {
		sourceNameI18nRaw = string(utils.MustMarshalJSON(g.SourceNameI18n))
	}

	gameInfo := &platformgame.GameInfo{
		Id:             g.ID,
		VenId:          g.ProviderID,
		CatId:          g.CategoryID,
		GroupId:        0,
		VenKey:         g.ProviderKey,
		Code:           g.GameCode,
		Name:           "",
		Image:          g.ImageURL,
		ChanId:         g.ChannelID,
		LoadType:       0,
		NameI18N:       nameI18nRaw,
		SourceId:       g.SourceID,
		SourceCode:     g.SourceGameCode,
		SourceNameI18N: sourceNameI18nRaw,
		SourceStatus:   int32(g.SourceStatus),
	}

	return &platformgame.GetGameResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data:    gameInfo,
	}, nil
}
