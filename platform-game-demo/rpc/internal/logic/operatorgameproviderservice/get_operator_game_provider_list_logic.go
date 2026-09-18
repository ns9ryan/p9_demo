package operatorgameproviderservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameProviderListLogic {
	return &GetOperatorGameProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏提供商列表
func (l *GetOperatorGameProviderListLogic) GetOperatorGameProviderList(in *platform_game.GetOperatorGameProviderListRequest) (*platform_game.GetOperatorGameProviderListResp, error) {
	l.Infof("[RPC GetOperatorGameProviderList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameProviderList] DAO Manager not available")
		return &platform_game.GetOperatorGameProviderListResp{}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	records, total, err := l.svcCtx.DAOManager.OperatorGameProvider.GetOperatorGameProviderList(
		l.ctx,
		in.GetOpCode(),
		in.GetProviderCode(),
		in.GetStatus(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameProviderList] query failed: %v", err)
		return &platform_game.GetOperatorGameProviderListResp{}, nil
	}

	l.Infof("[RPC GetOperatorGameProviderList] success: total=%d", total)
	return &platform_game.GetOperatorGameProviderListResp{
		Items:    logic.OperatorGameProviderModelToProtoList(records),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
