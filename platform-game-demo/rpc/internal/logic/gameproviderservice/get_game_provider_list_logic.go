package gameproviderservicelogic

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

type GetGameProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameProviderListLogic {
	return &GetGameProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏供应商列表
func (l *GetGameProviderListLogic) GetGameProviderList(in *platformgame.GetGameProviderListRequest) (*platformgame.GetGameProviderListResp, error) {
	l.Infof("[RPC GetGameProviderList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetGameProviderList] DAO Manager not available")
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	providers, total, err := l.svcCtx.DAOManager.GameProvider.GetGameProviderList(
		l.ctx,
		in.GetIsDeleted(),
		in.GetStatus(),
		in.GetProviderCode(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetGameProviderList] query failed: %v", err)
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	l.Infof("[RPC GetGameProviderList] success: total=%d", total)
	return &platformgame.GetGameProviderListResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Items:   logic.ProviderModelToProtoList(providers),
		Total:   int64(total),
	}, nil
}
