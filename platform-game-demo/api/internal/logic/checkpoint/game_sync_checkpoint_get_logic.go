// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package checkpoint

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/constant"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameSyncCheckpointGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameSyncCheckpointGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameSyncCheckpointGetLogic {
	return &GameSyncCheckpointGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameSyncCheckpointGetLogic) GameSyncCheckpointGet(req *types.GameSyncCheckpointGetReq) (resp *types.GameSyncCheckpointResp, err error) {
	l.Infof("[API GameSyncCheckpointGet] query checkpoint: sync_scope=%s", req.SyncScope)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Errorf("[API GameSyncCheckpointGet] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetGameSyncCheckpointServiceClient().GetGameSyncCheckpoint(l.ctx, &platformgame.GetGameSyncCheckpointRequest{
		Id:        req.ID,
		SyncScope: req.SyncScope,
	})
	if err != nil {
		l.Errorf("[API GameSyncCheckpointGet] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameSyncCheckpointGet] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	if grpcResp.Data == nil {
		l.Errorf("[API GameSyncCheckpointGet] gRPC data is nil")
		return nil, fmt.Errorf("gRPC data is nil")
	}

	resp = logic.CheckpointProtoToResponse(grpcResp.Data)

	l.Infof("[API GameSyncCheckpointGet] query success: sync_scope=%s", req.SyncScope)
	return resp, nil
}
