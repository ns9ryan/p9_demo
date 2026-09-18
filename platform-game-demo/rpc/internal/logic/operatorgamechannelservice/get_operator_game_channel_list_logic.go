package operatorgamechannelservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameChannelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameChannelListLogic {
	return &GetOperatorGameChannelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏渠道列表
func (l *GetOperatorGameChannelListLogic) GetOperatorGameChannelList(in *platform_game.GetOperatorGameChannelListRequest) (*platform_game.GetOperatorGameChannelListResp, error) {
	l.Infof("[RPC GetOperatorGameChannelList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameChannelList] DAO Manager not available")
		return &platform_game.GetOperatorGameChannelListResp{}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	records, total, err := l.svcCtx.DAOManager.OperatorGameChannel.GetOperatorGameChannelList(
		l.ctx,
		in.GetOpCode(),
		in.GetChannelCode(),
		in.GetStatus(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameChannelList] query failed: %v", err)
		return &platform_game.GetOperatorGameChannelListResp{}, nil
	}

	l.Infof("[RPC GetOperatorGameChannelList] success: total=%d", total)
	return &platform_game.GetOperatorGameChannelListResp{
		Items:    logic.OperatorGameChannelModelToProtoList(records),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
