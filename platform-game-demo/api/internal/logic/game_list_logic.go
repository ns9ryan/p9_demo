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

type GameListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameListLogic {
	return &GameListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameListLogic) GameList(req *types.GameListReq) (resp *types.GameListResp, err error) {
	l.Infof("[API GameList] received req: page=%d, page_size=%d, game_code=%s, name=%s, category_id=%d, provider_id=%d, channel_id=%d, status=%d, is_deleted=%d, sort_by='%s', sort_order='%s'",
		req.Page, req.PageSize, req.GameCode, req.Name, req.CategoryID, req.ProviderID, req.ChannelID, req.Status, req.IsDeleted, req.SortBy, req.SortOrder)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Errorf("[API GameList] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcReq := &platformgame.GetGameListRequest{
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		GameCode:   req.GameCode,
		Name:       req.Name,
		CategoryId: req.CategoryID,
		ProviderId: req.ProviderID,
		ChannelId:  req.ChannelID,
		Status:     int32(req.Status),
		IsDeleted:  int32(req.IsDeleted),
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetPlatformGameServiceClient().GetGameList(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameList] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameList] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	// 转换 proto 消息为 API 响应类型
	items := make([]types.GameResp, 0, len(grpcResp.Data))
	for _, data := range grpcResp.Data {
		l.Infof("[API GameList] rpc item: id=%d, cat_id=%d, ven_id=%d, chan_id=%d", data.Id, data.CatId, data.VenId, data.ChanId)
		items = append(items, types.GameResp{
			ID:             data.Id,
			GameCode:       data.Code,
			CategoryID:     data.CatId,
			ProviderID:     data.VenId,
			ChannelID:      data.ChanId,
			NameI18n:       data.NameI18N,
			Status:         int16(data.Status),
			ImageUrl:       data.Image,
			SourceID:       data.SourceId,
			SourceGameCode: data.SourceCode,
			SourceNameI18n: data.SourceNameI18N,
			SourceStatus:   int16(data.SourceStatus),
		})
		mapped := items[len(items)-1]
		l.Infof("[API GameList] mapped item: id=%d, category_id=%d, provider_id=%d, channel_id=%d", mapped.ID, mapped.CategoryID, mapped.ProviderID, mapped.ChannelID)
	}

	resp = &types.GameListResp{
		Items: items,
		Total: grpcResp.Total,
	}

	l.Infof("[API GameList] query success: total=%d, returned=%d", grpcResp.Total, len(items))
	return resp, nil
}
