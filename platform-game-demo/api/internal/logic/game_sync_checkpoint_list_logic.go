// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameSyncCheckpointListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameSyncCheckpointListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameSyncCheckpointListLogic {
	return &GameSyncCheckpointListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameSyncCheckpointListLogic) GameSyncCheckpointList(req *types.GameSyncCheckpointListReq) (resp *types.GameSyncCheckpointListResp, err error) {
	l.Infof("[API GameSyncCheckpointList] query checkpoint list: page=%d, page_size=%d, sync_scope=%s", req.Page, req.PageSize, req.SyncScope)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Errorf("[API GameSyncCheckpointList] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 15
	}
	if pageSize > 100 {
		pageSize = 100
	}

	grpcReq := &platformgame.GetGameSyncCheckpointListRequest{
		Page:      int32(page),
		PageSize:  int32(pageSize),
		SyncScope: req.SyncScope,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetPlatformGameServiceClient().GetGameSyncCheckpointList(l.ctx, grpcReq)
	if err != nil {
		l.Errorf("[API GameSyncCheckpointList] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameSyncCheckpointList] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	items := make([]types.GameSyncCheckpointResp, 0, len(grpcResp.Items))
	for _, item := range grpcResp.Items {
		items = append(items, types.GameSyncCheckpointResp{
			ID:              item.Id,
			SyncScope:       item.SyncScope,
			CheckpointValue: item.CheckpointValue,
			RemoteTotal:     item.RemoteTotal,
			LocalTotal:      item.LocalTotal,
			CreatedCount:    item.CreatedCount,
			UpdatedCount:    item.UpdatedCount,
			DeletedCount:    item.DeletedCount,
			LastSuccessAt:   item.LastSuccessAt,
			LastSyncAt:      item.LastSyncAt,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
		})
	}

	resp = &types.GameSyncCheckpointListResp{
		Items: items,
		Total: grpcResp.Total,
	}

	l.Infof("[API GameSyncCheckpointList] query success: total=%d, returned=%d", grpcResp.Total, len(items))
	return resp, nil
}
