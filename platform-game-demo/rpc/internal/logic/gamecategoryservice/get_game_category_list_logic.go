package gamecategoryservicelogic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameCategoryListLogic {
	return &GetGameCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏分类列表
func (l *GetGameCategoryListLogic) GetGameCategoryList(in *platform_game.GetGameCategoryListRequest) (*platform_game.GetGameCategoryListResp, error) {
	l.Infof("[RPC GetGameCategoryList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetGameCategoryList] DAO Manager not available")
		return &platform_game.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	categories, total, err := l.svcCtx.DAOManager.GameCategory.GetGameCategoryList(
		l.ctx,
		in.GetIsDeleted(),
		in.GetStatus(),
		in.GetCategoryCode(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetGameCategoryList] query failed: %v", err)
		return &platform_game.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	l.Infof("[RPC GetGameCategoryList] success: total=%d", total)
	return &platform_game.GetGameCategoryListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Items:    logic.GameCategoryModelToProtoList(categories),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
