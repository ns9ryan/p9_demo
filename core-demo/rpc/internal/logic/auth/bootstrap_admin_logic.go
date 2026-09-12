package auth

import (
	"context"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/bootstrap"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type BootstrapAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBootstrapAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BootstrapAdminLogic {
	return &BootstrapAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BootstrapAdminLogic) BootstrapAdmin(in *core.BootstrapAdminReq) (*core.UserPublic, error) {
	if l.svcCtx.Config.InitToken == "" || in.InitToken != l.svcCtx.Config.InitToken {
		return nil, xerr.RpcErr(xerr.Unauthorized(i18n.AuthInvalidInitToken))
	}
	u, err := bootstrap.CreatePlatformAdmin(l.ctx, l.svcCtx.Deps, bootstrap.CreateAdminReq{
		Username: in.Username, Password: in.Password, DisplayName: in.DisplayName,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	roles, err := l.svcCtx.Deps.RolesOfUser(l.ctx, u.ID)
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToUserPublic(service.PublicUser(u, roles)), nil
}
