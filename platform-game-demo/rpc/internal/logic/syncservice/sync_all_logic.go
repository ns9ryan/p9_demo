package syncservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"

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
