package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameSyncCheckpointLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameSyncCheckpointLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameSyncCheckpointLogic {
	return &GetGameSyncCheckpointLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取单个同步检查点
func (l *GetGameSyncCheckpointLogic) GetGameSyncCheckpoint(in *platformgame.GetGameSyncCheckpointRequest) (*platformgame.GetGameSyncCheckpointResp, error) {
	l.Infof("[RPC GetGameSyncCheckpoint] received request: sync_scope=%s", in.SyncScope)

	if l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameSyncCheckpoint] database not available")
		return &platformgame.GetGameSyncCheckpointResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	var checkpoint model.GameSyncCheckpoint
	query := l.svcCtx.DB.WithContext(l.ctx)

	// 如果指定了 sync_scope，按 sync_scope 查询；否则获取最近一条记录
	if in.SyncScope != "" {
		query = query.Where("sync_scope = ?", in.SyncScope)
	} else {
		query = query.Order("last_sync_at DESC")
	}

	if err := query.First(&checkpoint).Error; err != nil {
		l.Errorf("[RPC GetGameSyncCheckpoint] query failed: %v", err)
		return &platformgame.GetGameSyncCheckpointResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	info := &platformgame.GameSyncCheckpointInfo{
		Id:              checkpoint.ID,
		SyncScope:       checkpoint.SyncScope,
		CheckpointValue: checkpoint.CheckpointValue,
		RemoteTotal:     checkpoint.RemoteTotal,
		LocalTotal:      checkpoint.LocalTotal,
		CreatedCount:    checkpoint.CreatedCount,
		UpdatedCount:    checkpoint.UpdatedCount,
		DeletedCount:    checkpoint.DeletedCount,
		LastSyncAt:      checkpoint.LastSyncAt.Unix(),
		LastSuccessAt:   checkpoint.LastSuccessAt.Unix(),
		LastError:       derefString(checkpoint.LastErrorMessage),
		CreatedAt:       checkpoint.CreatedAt.Unix(),
		UpdatedAt:       checkpoint.UpdatedAt.Unix(),
	}

	l.Infof("[RPC GetGameSyncCheckpoint] query success: sync_scope=%s", checkpoint.SyncScope)
	return &platformgame.GetGameSyncCheckpointResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    info,
	}, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
