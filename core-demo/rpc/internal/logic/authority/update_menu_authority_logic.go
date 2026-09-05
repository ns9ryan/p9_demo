package authority

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuAuthorityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuAuthorityLogic {
	return &UpdateMenuAuthorityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuAuthorityLogic) UpdateMenuAuthority(in *core.MenuAuthReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.UpdateMenuAuthority(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.MenuAuthReq{RoleID: in.RoleId, MenuIDs: in.MenuIds})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
