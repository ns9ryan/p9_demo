package user

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserRolesLogic {
	return &BindUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindUserRolesLogic) BindUserRoles(req *types.BindRolesReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.BindUserRoles(l.ctx, &coreclient.BindRolesReq{UserId: req.UserId, RoleIds: req.RoleIds})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Result: "success"}, nil
}
