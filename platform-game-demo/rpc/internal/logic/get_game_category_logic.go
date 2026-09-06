package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/pkg/game"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

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
func (l *GetGameCategoryLogic) GetGameCategory(in *platformgame.GetGameCategoryRequest) (*platformgame.GetGameCategoryListResp, error) {
	l.Infof("[RPC GetGameCategory] received request: id=%d", in.Id)

	if l.svcCtx == nil {
		l.Error("[RPC GetGameCategory] ServiceContext is nil")
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeSuccess,
			Message: "ServiceContext is nil",
			Data:    nil,
		}, nil
	}

	if l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameCategory] Database not available (nil). Config: Driver=%s, DSN=%s",
			l.svcCtx.Config.Database.Driver, l.svcCtx.Config.Database.DSN)
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available (nil)",
			Data:    nil,
		}, nil
	}

	l.Infof("[RPC GetGameCategory] Database is available")

	// Call pkg/game to get category
	params := &game.GameCategoryGetParams{
		ID: in.Id,
	}
	category, err := game.GameCategoryGetDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("[RPC GetGameCategory] failed to get category: %v", err)
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get category: " + err.Error(),
			Data:    nil,
		}, nil
	}

	if category == nil {
		l.Infof("[RPC GetGameCategory] category not found: id=%d", in.Id)
		return &platformgame.GetGameCategoryListResp{
			Code:    constant.CodeNotFound,
			Message: "category not found",
			Data:    nil,
		}, nil
	}

	l.Infof("[RPC GetGameCategory] found category: id=%d, code=%s", category.ID, category.CategoryCode)

	// Convert to proto message
	item := &platformgame.CategoryInfo{
		Id:             category.ID,
		Code:           category.CategoryCode,
		Name:           category.CategoryCode,
		Status:         int32(category.Status),
		NameI18N:       string(utils.MustMarshalJSON(category.NameI18n)),
		SourceNameI18N: string(utils.MustMarshalJSON(category.SourceNameI18n)),
	}

	return &platformgame.GetGameCategoryListResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    []*platformgame.CategoryInfo{item},
	}, nil
}
