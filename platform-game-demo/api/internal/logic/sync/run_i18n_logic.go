// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package sync

import (
	"context"

	"oa.98ent.com/p9/platform-game/api/internal/catalog"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunI18nLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRunI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunI18nLogic {
	return &RunI18nLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunI18nLogic) RunI18n(req *types.SyncRunReq) (resp *types.SyncRunResp, err error) {
	r, err := l.svcCtx.GrpcClient.GetGameSyncCheckpointServiceClient().
		GetI18NNameMap(context.Background(), &platform_game.GetI18NNameMapRequest{})
	if err == nil {
		l.svcCtx.Core.RegisterCatalog(context.Background(), catalog.AppendGameI18nItems(r.Data))
		l.Infof("[API GameSyncCheckpointGet] i18n name map registered successfully")
	} else {
		l.Errorf("[API GameSyncCheckpointGet] failed to get i18n name map: %v", err)
	}
	return &types.SyncRunResp{}, nil
}
