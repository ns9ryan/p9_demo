package operatorgamecategoryservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameCategoryListLogic {
	return &GetOperatorGameCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏分类列表
func (l *GetOperatorGameCategoryListLogic) GetOperatorGameCategoryList(in *platform_game.GetOperatorGameCategoryListRequest) (*platform_game.GetOperatorGameCategoryListResp, error) {
	l.Infof("[RPC GetOperatorGameCategoryList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameCategoryList] DAO Manager not available")
		return &platform_game.GetOperatorGameCategoryListResp{}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	records, total, err := l.svcCtx.DAOManager.OperatorGameCategory.GetOperatorGameCategoryList(
		l.ctx,
		in.GetOpCode(),
		in.GetCategoryCode(),
		in.GetStatus(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameCategoryList] query failed: %v", err)
		return &platform_game.GetOperatorGameCategoryListResp{}, nil
	}

	l.Infof("[RPC GetOperatorGameCategoryList] success: total=%d", total)
	return &platform_game.GetOperatorGameCategoryListResp{
		Items:    logic.OperatorGameCategoryModelToProtoList(records),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
