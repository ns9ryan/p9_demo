package authority

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuAuthorityLogic {
	return &GetMenuAuthorityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuAuthorityLogic) GetMenuAuthority(req *types.RoleIdReq) (resp *types.MenuAuthResp, err error) {
	out, err := l.svcCtx.Core.GetMenuAuthority(l.ctx, &coreclient.RoleIdReq{RoleId: req.RoleId, Id: req.Id})
	if err != nil {
		return nil, err
	}
	ids := out.GetMenuIds()
	if ids == nil {
		ids = []int64{}
	}
	return &types.MenuAuthResp{MenuIds: ids}, nil
}
