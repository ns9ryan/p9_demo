package gameservicelogic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
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
	l.Infof("[RPC GetGameList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetGameList] DAO Manager not available")
		return &platformgame.GetGameListResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	games, total, err := l.svcCtx.DAOManager.Game.GetGameList(
		l.ctx,
		in.GetIsDeleted(),
		in.GetName(),
		in.GetStatus(),
		in.GetGameCode(),
		int64(in.GetProviderId()),
		int64(in.GetCategoryId()),
		int64(in.GetChannelId()),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetGameList] query failed: %v", err)
		return &platformgame.GetGameListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	var gameProtoList []*platformgame.GameInfo
	for _, gameRecord := range games {
		ext, err := GetGameExtInfo(l.ctx, l.svcCtx, gameRecord)
		if err != nil {
			l.Errorf("[RPC GetGameList] query extended game info failed: %v", err)
			return &platformgame.GetGameListResp{
				Code:    constant.CodeInternalError,
				Message: fmt.Sprintf("query extended game info failed: %v", err),
			}, nil
		}
		gameProtoList = append(gameProtoList, logic.GameModelToProto(gameRecord, ext))
	}

	l.Infof("[RPC GetGameList] success: total=%d", total)
	return &platformgame.GetGameListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Items:    gameProtoList,
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
