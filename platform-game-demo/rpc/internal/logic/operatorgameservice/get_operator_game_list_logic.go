package operatorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameListLogic {
	return &GetOperatorGameListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏列表
func (l *GetOperatorGameListLogic) GetOperatorGameList(in *platform_game.GetOperatorGameListRequest) (*platform_game.GetOperatorGameListResp, error) {
	l.Infof("[RPC GetOperatorGameList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameList] DAO Manager not available")
		return &platform_game.GetOperatorGameListResp{}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	records, total, err := l.svcCtx.DAOManager.OperatorGame.GetOperatorGameList(
		l.ctx,
		in.GetOpCode(),
		in.GetGameCode(),
		in.GetStatus(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameList] query failed: %v", err)
		return &platform_game.GetOperatorGameListResp{}, nil
	}

	l.Infof("[RPC GetOperatorGameList] success: total=%d", total)
	return &platform_game.GetOperatorGameListResp{
		Items:    logic.OperatorGameModelToProtoList(records),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
