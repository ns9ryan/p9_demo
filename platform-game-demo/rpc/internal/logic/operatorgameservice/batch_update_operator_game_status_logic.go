package operatorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateOperatorGameStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUpdateOperatorGameStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameStatusLogic {
	return &BatchUpdateOperatorGameStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量修改分站游戏状态
func (l *BatchUpdateOperatorGameStatusLogic) BatchUpdateOperatorGameStatus(in *platform_game.BatchUpdateOperatorGameStatusRequest) (*platform_game.BatchUpdateOperatorGameStatusResp, error) {
	l.Infof("[RPC BatchUpdateOperatorGameStatus] received req: %d ids", len(in.Ids))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchUpdateOperatorGameStatus] DAO Manager not available")
		return &platform_game.BatchUpdateOperatorGameStatusResp{}, nil
	}

	if len(in.Ids) == 0 {
		return &platform_game.BatchUpdateOperatorGameStatusResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	affected, err := l.svcCtx.DAOManager.OperatorGame.BatchUpdateOperatorGameStatus(l.ctx, in.Ids, int16(in.Status))
	if err != nil {
		l.Errorf("[RPC BatchUpdateOperatorGameStatus] batch update failed: %v", err)
		return &platform_game.BatchUpdateOperatorGameStatusResp{
			Total:   int64(len(in.Ids)),
			Success: 0,
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	l.Infof("[RPC BatchUpdateOperatorGameStatus] success: updated %d records", affected)
	return &platform_game.BatchUpdateOperatorGameStatusResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids) - affected),
	}, nil
}
