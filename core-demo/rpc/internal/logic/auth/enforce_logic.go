package auth

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnforceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEnforceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnforceLogic {
	return &EnforceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EnforceLogic) Enforce(in *core.EnforceReq) (*core.EnforceResp, error) {
	claims := ctxdata.ClaimsFromCtx(l.ctx)
	if claims == nil {
		claims = &ctxdata.Claims{RoleCodes: in.RoleCodes, OperatorID: in.OperatorId}
	}
	ok, err := l.svcCtx.Deps.Enforce(l.ctx, claims, in.Path, in.Method)
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.EnforceResp{Allowed: ok}, nil
}
