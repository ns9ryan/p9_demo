package operatorgamecategoryservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateOperatorGameCategoryStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUpdateOperatorGameCategoryStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameCategoryStatusLogic {
	return &BatchUpdateOperatorGameCategoryStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量修改分站游戏分类状态
func (l *BatchUpdateOperatorGameCategoryStatusLogic) BatchUpdateOperatorGameCategoryStatus(in *platform_game.BatchUpdateOperatorGameCategoryStatusRequest) (*platform_game.BatchUpdateOperatorGameCategoryStatusResp, error) {
	l.Infof("[RPC BatchUpdateOperatorGameCategoryStatus] received req: %d ids", len(in.Ids))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchUpdateOperatorGameCategoryStatus] DAO Manager not available")
		return &platform_game.BatchUpdateOperatorGameCategoryStatusResp{}, nil
	}

	if len(in.Ids) == 0 {
		return &platform_game.BatchUpdateOperatorGameCategoryStatusResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	affected, err := l.svcCtx.DAOManager.OperatorGameCategory.BatchUpdateOperatorGameCategoryStatus(l.ctx, in.Ids, int16(in.Status))
	if err != nil {
		l.Errorf("[RPC BatchUpdateOperatorGameCategoryStatus] batch update failed: %v", err)
		return &platform_game.BatchUpdateOperatorGameCategoryStatusResp{
			Total:   int64(len(in.Ids)),
			Success: 0,
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	l.Infof("[RPC BatchUpdateOperatorGameCategoryStatus] success: updated %d records", affected)
	return &platform_game.BatchUpdateOperatorGameCategoryStatusResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids) - affected),
	}, nil
}
