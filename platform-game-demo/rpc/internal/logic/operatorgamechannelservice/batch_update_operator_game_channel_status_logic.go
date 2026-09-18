package operatorgamechannelservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateOperatorGameChannelStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUpdateOperatorGameChannelStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateOperatorGameChannelStatusLogic {
	return &BatchUpdateOperatorGameChannelStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量修改分站游戏渠道状态
func (l *BatchUpdateOperatorGameChannelStatusLogic) BatchUpdateOperatorGameChannelStatus(in *platform_game.BatchUpdateOperatorGameChannelStatusRequest) (*platform_game.BatchUpdateOperatorGameChannelStatusResp, error) {
	l.Infof("[RPC BatchUpdateOperatorGameChannelStatus] received req: %d ids", len(in.Ids))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchUpdateOperatorGameChannelStatus] DAO Manager not available")
		return &platform_game.BatchUpdateOperatorGameChannelStatusResp{}, nil
	}

	if len(in.Ids) == 0 {
		return &platform_game.BatchUpdateOperatorGameChannelStatusResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	affected, err := l.svcCtx.DAOManager.OperatorGameChannel.BatchUpdateOperatorGameChannelStatus(l.ctx, in.Ids, int16(in.Status))
	if err != nil {
		l.Errorf("[RPC BatchUpdateOperatorGameChannelStatus] batch update failed: %v", err)
		return &platform_game.BatchUpdateOperatorGameChannelStatusResp{
			Total:   int64(len(in.Ids)),
			Success: 0,
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	l.Infof("[RPC BatchUpdateOperatorGameChannelStatus] success: updated %d records", affected)
	return &platform_game.BatchUpdateOperatorGameChannelStatusResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids) - affected),
	}, nil
}
