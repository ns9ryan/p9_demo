package operatorgamecategoryservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteOperatorGameCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteOperatorGameCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperatorGameCategoryLogic {
	return &BatchDeleteOperatorGameCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除分站游戏分类
func (l *BatchDeleteOperatorGameCategoryLogic) BatchDeleteOperatorGameCategory(in *platform_game.BatchDeleteOperatorGameCategoryRequest) (*platform_game.BatchDeleteOperatorGameCategoryResp, error) {
	affected, err := l.svcCtx.DAOManager.OperatorGameCategory.BatchDeleteOperatorGameCategory(l.ctx, in.Ids)
	if err != nil {
		return &platform_game.BatchDeleteOperatorGameCategoryResp{
			Total:   int64(len(in.Ids)),
			Success: int64(0),
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	return &platform_game.BatchDeleteOperatorGameCategoryResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids)) - int64(affected),
	}, nil
}
