package gamecategoryservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameCategoryLogic {
	return &GetGameCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取单个游戏分类
func (l *GetGameCategoryLogic) GetGameCategory(in *platform_game.GetGameCategoryRequest) (*platform_game.GetGameCategoryResp, error) {
	l.Infof("[RPC GetGameCategory] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameCategory] Database not available")
		return &platform_game.GetGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	category := &ent.GameCategory{}
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(category).Error; err != nil {
		l.Errorf("[RPC GetGameCategory] query failed: %v", err)
		return &platform_game.GetGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	return &platform_game.GetGameCategoryResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.CategoryModelToProto(category),
	}, nil
}
