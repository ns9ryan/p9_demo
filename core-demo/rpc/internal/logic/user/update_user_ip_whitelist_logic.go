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

type UpdateUserIpWhitelistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserIpWhitelistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserIpWhitelistLogic {
	return &UpdateUserIpWhitelistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserIpWhitelistLogic) UpdateUserIpWhitelist(in *core.UpdateUserIpWhitelistReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.UpdateUserIpWhitelist(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.UpdateUserIpWhitelistReq{
		ID: in.Id, IPWhitelistEnabled: int16(in.IpWhitelistEnabled), IPWhitelist: in.IpWhitelist,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
