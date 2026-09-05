package authority

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuAuthorityLogic {
	return &UpdateMenuAuthorityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuAuthorityLogic) UpdateMenuAuthority(req *types.MenuAuthReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.UpdateMenuAuthority(l.ctx, &coreclient.MenuAuthReq{RoleId: req.RoleId, MenuIds: req.MenuIds})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Result: "success"}, nil
}
