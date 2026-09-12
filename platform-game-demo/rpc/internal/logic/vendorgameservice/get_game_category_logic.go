package vendorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"

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

// 获取游戏分类
func (l *GetGameCategoryLogic) GetGameCategory(in *vendors.Empty) (*vendors.GetGameCategoryResponse, error) {
	// todo: add your logic here and delete this line

	return &vendors.GetGameCategoryResponse{}, nil
}
