package operatorgameproviderservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteOperatorGameProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteOperatorGameProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperatorGameProviderLogic {
	return &BatchDeleteOperatorGameProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除分站游戏提供商
func (l *BatchDeleteOperatorGameProviderLogic) BatchDeleteOperatorGameProvider(in *platform_game.BatchDeleteOperatorGameProviderRequest) (*platform_game.BatchDeleteOperatorGameProviderResp, error) {
	affected, err := l.svcCtx.DAOManager.OperatorGameProvider.BatchDeleteOperatorGameProvider(l.ctx, in.Ids)
	if err != nil {
		return &platform_game.BatchDeleteOperatorGameProviderResp{
			Total:   int64(len(in.Ids)),
			Success: int64(0),
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	return &platform_game.BatchDeleteOperatorGameProviderResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids)) - int64(affected),
	}, nil
}
