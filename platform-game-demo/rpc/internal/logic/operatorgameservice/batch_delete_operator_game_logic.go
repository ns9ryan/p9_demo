package operatorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteOperatorGameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteOperatorGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperatorGameLogic {
	return &BatchDeleteOperatorGameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除分站游戏
func (l *BatchDeleteOperatorGameLogic) BatchDeleteOperatorGame(in *platform_game.BatchDeleteOperatorGameRequest) (*platform_game.BatchDeleteOperatorGameResp, error) {
	affected, err := l.svcCtx.DAOManager.OperatorGame.BatchDeleteOperatorGame(l.ctx, in.Ids)
	if err != nil {
		return &platform_game.BatchDeleteOperatorGameResp{
			Total:   int64(len(in.Ids)),
			Success: int64(0),
			Failed:  int64(len(in.Ids)),
		}, nil
	}

	return &platform_game.BatchDeleteOperatorGameResp{
		Total:   int64(len(in.Ids)),
		Success: int64(affected),
		Failed:  int64(len(in.Ids)) - int64(affected),
	}, nil
}
