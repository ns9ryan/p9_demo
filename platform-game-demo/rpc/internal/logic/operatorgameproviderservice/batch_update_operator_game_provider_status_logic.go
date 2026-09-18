package operatorgameproviderservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateOperatorGameProviderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUpdateOperatorGameProviderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameProviderStatusLogic {
	return &BatchUpdateOperatorGameProviderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量修改分站游戏提供商状态
func (l *BatchUpdateOperatorGameProviderStatusLogic) BatchUpdateOperatorGameProviderStatus(in *platform_game.BatchUpdateOperatorGameProviderStatusRequest) (*platform_game.BatchUpdateOperatorGameProviderStatusResp, error) {
	l.Infof("[RPC BatchUpdateOperatorGameProviderStatus] received req: %d ids", len(in.Ids))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchUpdateOperatorGameProviderStatus] DAO Manager not available")
		return &platform_game.BatchUpdateOperatorGameProviderStatusResp{}, nil
	}

	if len(in.Ids) == 0 {
		return &platform_game.BatchUpdateOperatorGameProviderStatusResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	affected, err := l.svcCtx.DAOManager.OperatorGameProvider.BatchUpdateOperatorGameProviderStatus(l.ctx, in.Ids, int16(in.Status))
	if err != nil {
		l.Errorf("[RPC BatchUpdateOperatorGameProviderStatus] batch update failed: %v", err)
		return &platform_game.BatchUpdateOperatorGameProviderStatusResp{
			Total:   int64(len(in.Ids)),
			Success: 0,
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	l.Infof("[RPC BatchUpdateOperatorGameProviderStatus] success: updated %d records", affected)
	return &platform_game.BatchUpdateOperatorGameProviderStatusResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids) - affected),
	}, nil
}
