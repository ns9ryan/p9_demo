// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package game

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

type GameGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameGetLogic {
	return &GameGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameGetLogic) GameGet(req *types.GameGetReq) (resp *types.GameResp, err error) {
	l.Infof("[API GameGet] query game: id=%d", req.ID)

	if l.svcCtx == nil || l.svcCtx.GrpcClient == nil {
		l.Errorf("[API GameGet] gRPC client not available")
		return nil, fmt.Errorf("gRPC client not available")
	}

	grpcResp, err := l.svcCtx.GrpcClient.GetGameServiceClient().GetGame(l.ctx, &platformgame.GetGameRequest{
		Id: req.ID,
	})
	if err != nil {
		l.Errorf("[API GameGet] gRPC call failed: %v", err)
		return nil, fmt.Errorf("gRPC call failed: %s", err.Error())
	}

	if grpcResp.Code != constant.CodeSuccess {
		l.Errorf("[API GameGet] gRPC error: Code=%d, Message=%s", grpcResp.Code, grpcResp.Message)
		return nil, fmt.Errorf("gRPC error: %s", grpcResp.Message)
	}

	if grpcResp.Data == nil {
		l.Errorf("[API GameGet] gRPC data is nil")
		return nil, fmt.Errorf("gRPC data is nil")
	}

	// 转换 proto 消息为 API 响应类型
	resp = logic.GameProtoToResponse(grpcResp.Data)

	l.Infof("[API GameGet] query success: id=%d, source_id=%d", req.ID, grpcResp.Data.SourceId)
	return resp, nil
}
