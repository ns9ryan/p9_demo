package logic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/pkg/game"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameSyncCheckpointListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameSyncCheckpointListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameSyncCheckpointListLogic {
	return &GetGameSyncCheckpointListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取同步检查点列表
func (l *GetGameSyncCheckpointListLogic) GetGameSyncCheckpointList(in *platformgame.GetGameSyncCheckpointListRequest) (*platformgame.GetGameSyncCheckpointListResp, error) {
	l.Infof("[RPC GetGameSyncCheckpointList] received request: page=%d, page_size=%d, sync_scope=%q, start_time=%d, end_time=%d", in.Page, in.PageSize, in.SyncScope, in.StartTime, in.EndTime)

	if l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameSyncCheckpointList] database not available")
		return &platformgame.GetGameSyncCheckpointListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 15
	}
	if pageSize > 100 {
		pageSize = 100
	}

	params := &game.GameSyncCheckpointListParams{
		Page:      int64(page),
		PageSize:  int64(pageSize),
		SyncScope: in.SyncScope,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
	}

	var startAt string
	if in.StartTime > 0 {
		startAt = time.Unix(in.StartTime, 0).Format(time.RFC3339)
	}
	var endAt string
	if in.EndTime > 0 {
		endAt = time.Unix(in.EndTime, 0).Format(time.RFC3339)
	}
	l.Infof("[RPC GetGameSyncCheckpointList] normalized params: page=%d, page_size=%d, sync_scope=%q, start_at=%q, end_at=%q", page, pageSize, in.SyncScope, startAt, endAt)

	result, err := game.GameSyncCheckpointListDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("[RPC GetGameSyncCheckpointList] query failed: %v", err)
		return &platformgame.GetGameSyncCheckpointListResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	items := make([]*platformgame.GameSyncCheckpointInfo, 0, len(result.Checkpoints))
	for _, cp := range result.Checkpoints {
		items = append(items, &platformgame.GameSyncCheckpointInfo{
			Id:              cp.ID,
			SyncScope:       cp.SyncScope,
			CheckpointValue: cp.CheckpointValue,
			RemoteTotal:     cp.RemoteTotal,
			LocalTotal:      cp.LocalTotal,
			CreatedCount:    cp.CreatedCount,
			UpdatedCount:    cp.UpdatedCount,
			DeletedCount:    cp.DeletedCount,
			LastSyncAt:      cp.LastSyncAt.Unix(),
			LastSuccessAt:   cp.LastSuccessAt.Unix(),
			LastError:       derefString(cp.LastErrorMessage),
			CreatedAt:       cp.CreatedAt.Unix(),
			UpdatedAt:       cp.UpdatedAt.Unix(),
		})
	}

	l.Infof("[RPC GetGameSyncCheckpointList] query success: total=%d, returned=%d", result.Total, len(items))
	return &platformgame.GetGameSyncCheckpointListResp{
		Code:     constant.CodeSuccess,
		Message:  "ok",
		Items:    items,
		Total:    result.Total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
