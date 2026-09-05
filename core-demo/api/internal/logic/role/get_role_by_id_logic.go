package role

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleByIdLogic {
	return &GetRoleByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleByIdLogic) GetRoleById(req *types.IDQuery) (resp *types.RoleInfo, err error) {
	out, err := l.svcCtx.Core.GetRoleById(l.ctx, &coreclient.IDReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return convert.RoleInfo(l.ctx, out), nil
}
