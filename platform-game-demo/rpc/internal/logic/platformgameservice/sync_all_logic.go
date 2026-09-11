package platformgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	sync "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncAllLogic {
	return &SyncAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 全量同步（同步所有对象类型）
func (l *SyncAllLogic) SyncAll(in *sync.SyncAllRequest) (*sync.SyncRunResp, error) {
	// todo: add your logic here and delete this line

	return &sync.SyncRunResp{}, nil
}
