package gamechannelservicelogic

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

type GetGameChannelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameChannelListLogic {
	return &GetGameChannelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏渠道列表
func (l *GetGameChannelListLogic) GetGameChannelList(in *platformgame.GetGameChannelListRequest) (*platformgame.GetGameChannelListResp, error) {
	l.Infof("[RPC GetGameChannelList] received req: page=%d, page_size=%d", in.Page, in.PageSize)

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Error("[RPC GetGameChannelList] DAO Manager not available")
		return &platformgame.GetGameChannelListResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	channels, total, err := l.svcCtx.DAOManager.GameChannel.GetGameChannelList(
		l.ctx,
		in.GetIsDeleted(),
		in.GetStatus(),
		in.GetChannelCode(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetGameChannelList] query failed: %v", err)
		return &platformgame.GetGameChannelListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	items := logic.ChannelModelToProtoList(channels)
	l.Infof("[RPC GetGameChannelList] query result: returned=%d items, total=%d", len(items), total)

	return &platformgame.GetGameChannelListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Items:    items,
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
