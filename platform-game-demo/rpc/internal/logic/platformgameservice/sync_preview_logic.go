package platformgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	sync "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncPreviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncPreviewLogic {
	return &SyncPreviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步预检查（查看差异，不修改数据）
func (l *SyncPreviewLogic) SyncPreview(in *sync.SyncPreviewRequest) (*sync.SyncPreviewResp, error) {
	// todo: add your logic here and delete this line

	return &sync.SyncPreviewResp{}, nil
}
