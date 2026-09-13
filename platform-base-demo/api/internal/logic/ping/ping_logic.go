// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package ping

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/api/internal/types"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/ping"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Ping 检查服务是否正常
func (l *PingLogic) Ping() (resp *types.PingResponse, err error) {
	// 调用Ping RPC
	_, err = l.svcCtx.PingRpc.Ping(
		l.ctx,
		&ping.PingRequest{},
	)
	if err != nil {
		return nil, err
	}

	return &types.PingResponse{}, nil
}
