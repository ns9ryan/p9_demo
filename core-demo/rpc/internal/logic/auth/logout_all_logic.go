package auth

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutAllLogic {
	return &LogoutAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutAllLogic) LogoutAll(in *core.Empty) (*core.Empty, error) {
	if err := l.svcCtx.Deps.LogoutAll(l.ctx); err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
