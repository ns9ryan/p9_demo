package user

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserRolesLogic {
	return &BindUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindUserRolesLogic) BindUserRoles(in *core.BindRolesReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.BindUserRoles(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.BindRolesReq{UserID: in.UserId, RoleIDs: in.RoleIds})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
