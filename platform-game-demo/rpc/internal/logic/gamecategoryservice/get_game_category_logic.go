package gamecategoryservicelogic

import (
	"context"
	"fmt"

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
	l.Infof("[RPC GetGameCategory] received request: id=%d", in.GetId())

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetGameCategory] DAO Manager not available")
		return &platform_game.GetGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	category, err := l.svcCtx.DAOManager.GameCategory.GetGameCategoryByID(l.ctx, in.GetId())
	if err != nil {
		l.Errorf("[RPC GetGameCategory] query failed: %v", err)
		return &platform_game.GetGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("failed to get game category: %v", err),
		}, nil
	}

	l.Infof("[RPC GetGameCategory] success: id=%d", category.ID)
	return &platform_game.GetGameCategoryResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.GameCategoryModelToProto(category),
	}, nil
}
