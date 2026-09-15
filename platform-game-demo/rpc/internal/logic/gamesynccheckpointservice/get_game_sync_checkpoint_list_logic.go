package gamesynccheckpointservicelogic

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent/gamesynccheckpoint"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
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
	l.Infof("[RPC GetGameSyncCheckpointList] received req: page=%d, page_size=%d", in.Page, in.PageSize)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Error("[RPC GetGameSyncCheckpointList] database not available")
		return &platformgame.GetGameSyncCheckpointListResp{
			Code:    constant.CodeInternalError,
			Message: "database not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))

	query := l.svcCtx.DB.GameSyncCheckpoint.Query()

	if in.GetSyncScope() != "" {
		query = query.Where(gamesynccheckpoint.SyncScopeEQ(in.GetSyncScope()))
	}

	if in.GetStartTime() > 0 {
		startTime := time.Unix(in.GetStartTime(), 0)
		query = query.Where(gamesynccheckpoint.CreatedAtGTE(startTime))
	}

	if in.GetEndTime() > 0 {
		endTime := time.Unix(in.GetEndTime(), 0)
		query = query.Where(gamesynccheckpoint.CreatedAtLTE(endTime))
	}

	sortBy := "created_at"
	sortOrder := "desc"

	switch sortBy {
	case "id", "created_at", "last_sync_at":
	default:
		sortBy = "created_at"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		l.Errorf("[RPC GetGameSyncCheckpointList] count failed: %v", err)
		return &platformgame.GetGameSyncCheckpointListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("count failed: %v", err),
		}, nil
	}
	// 计算分页偏移量
	offset := (page - 1) * pageSize

	// 获取当前页分类数据
	checkpoints, err := query.
		Order(
			gamesynccheckpoint.ByCreatedAt(), // 按创建时间升序
			gamesynccheckpoint.ByID(),        // 创建时间相同时按ID升序
		).
		Offset(int(offset)).
		Limit(int(pageSize)).
		All(l.ctx)
	if err != nil {
		l.Errorf("[RPC GetGameSyncCheckpointList] query failed: %v", err)
		return &platformgame.GetGameSyncCheckpointListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	l.Infof("[RPC GetGameSyncCheckpointList] success: total=%d", total)
	return &platformgame.GetGameSyncCheckpointListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Items:    logic.CheckpointModelToProtoList(checkpoints),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
