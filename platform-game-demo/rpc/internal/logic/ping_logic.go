package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Ping 健康检查
func (l *PingLogic) Ping(in *platformgame.PingRequest) (*platformgame.PingResponse, error) {
	l.Infof("Ping RPC called")
	return &platformgame.PingResponse{
		Message: "pong",
	}, nil
}
