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

type GetGameListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameListLogic {
	return &GetGameListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏列表
func (l *GetGameListLogic) GetGameList(in *platformgame.GetGameListRequest) (*platformgame.GetGameListResp, error) {
	l.Infof("GetGameList called with page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx.DB == nil {
		l.Errorf("Database not available")
		return &platformgame.GetGameListResp{
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
	if pageSize > 100 {
		pageSize = 100
	}

	params := &game.GameListParams{
		Page:       page,
		PageSize:   pageSize,
		GameCode:   in.GameCode,
		Name:       in.Name,
		CategoryID: in.CategoryId,
		ProviderID: in.ProviderId,
		ChannelID:  in.ChannelId,
		Status:     int16(in.Status),
		IsDeleted:  int8(in.IsDeleted),
		SortBy:     in.GetSortBy(),
		SortOrder:  in.GetSortOrder(),
	}

	// 调用数据库查询函数
	result, err := game.GameListDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("GameListDB failed: %v", err)
		return &platformgame.GetGameListResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to fetch games: " + err.Error(),
		}, nil
	}

	// 转换结果为protobuf 消息
	games := make([]*platformgame.GameInfo, 0, len(result.Games))
	for _, item := range result.Games {
		if item.Game == nil {
			continue
		}

		g := item.Game
		l.Infof("[RPC GetGameList] DB game row: game_id=%d, category_id=%d, provider_id=%d, channel_id=%d, source_id=%d, source_game_code=%s, source_status=%d", g.ID, g.CategoryID, g.ProviderID, g.ChannelID, g.SourceID, g.SourceGameCode, g.SourceStatus)

		games = append(games, &platformgame.GameInfo{
			Id:             g.ID,
			VenId:          g.ProviderID,
			CatId:          g.CategoryID,
			GroupId:        0,
			VenKey:         g.ProviderKey,
			Code:           g.GameCode,
			Name:           g.GameCode,
			Image:          g.ImageURL,
			ChanId:         g.ChannelID,
			LoadType:       1,
			Status:         int32(g.Status),
			NameI18N:       string(utils.MustMarshalJSON(g.NameI18n)),
			SourceId:       g.SourceID,
			SourceCode:     g.SourceGameCode,
			SourceNameI18N: string(utils.MustMarshalJSON(g.SourceNameI18n)),
			SourceStatus:   int32(g.SourceStatus),
		})
		mapped := games[len(games)-1]
		l.Infof("[RPC GetGameList] mapped GameInfo: id=%d, cat_id=%d, ven_id=%d, chan_id=%d, source_id=%d, source_code=%s, source_status=%d", mapped.Id, mapped.CatId, mapped.VenId, mapped.ChanId, mapped.SourceId, mapped.SourceCode, mapped.SourceStatus)
	}

	l.Infof("GetGameList succeeded: %d games found, total=%d", len(games), result.Total)
	return &platformgame.GetGameListResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data:    games,
	}, nil
}
