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

type UpdateApiAuthorityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApiAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiAuthorityLogic {
	return &UpdateApiAuthorityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateApiAuthorityLogic) UpdateApiAuthority(in *core.ApiAuthReq) (*core.Empty, error) {
	data := make([]service.APIAuthItem, 0, len(in.Data))
	for _, it := range in.Data {
		data = append(data, service.APIAuthItem{Path: it.Path, Method: it.Method})
	}
	err := l.svcCtx.Deps.UpdateAPIAuthority(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.APIAuthReq{RoleID: in.RoleId, Data: data})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
