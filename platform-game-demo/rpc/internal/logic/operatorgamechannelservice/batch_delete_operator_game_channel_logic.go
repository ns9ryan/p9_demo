package operatorgamechannelservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteOperatorGameChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteOperatorGameChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperatorGameChannelLogic {
	return &BatchDeleteOperatorGameChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除分站游戏渠道
func (l *BatchDeleteOperatorGameChannelLogic) BatchDeleteOperatorGameChannel(in *platform_game.BatchDeleteOperatorGameChannelRequest) (*platform_game.BatchDeleteOperatorGameChannelResp, error) {
	affected, err := l.svcCtx.DAOManager.OperatorGameChannel.BatchDeleteOperatorGameChannel(l.ctx, in.Ids)
	if err != nil {
		return &platform_game.BatchDeleteOperatorGameChannelResp{
			Total:   int64(len(in.Ids)),
			Success: int64(0),
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	return &platform_game.BatchDeleteOperatorGameChannelResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids)) - int64(affected),
	}, nil
}
