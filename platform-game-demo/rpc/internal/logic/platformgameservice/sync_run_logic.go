package platformgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	sync "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncRunLogic {
	return &SyncRunLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步执行（执行同步操作）
func (l *SyncRunLogic) SyncRun(in *sync.SyncRunRequest) (*sync.SyncRunResp, error) {
	// todo: add your logic here and delete this line

	return &sync.SyncRunResp{}, nil
}
