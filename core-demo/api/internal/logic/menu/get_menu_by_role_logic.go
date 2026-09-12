package menu

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuByRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuByRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuByRoleLogic {
	return &GetMenuByRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuByRoleLogic) GetMenuByRole() (resp []types.MenuNode, err error) {
	out, err := l.svcCtx.Core.GetMenuByRole(l.ctx, &coreclient.Empty{})
	if err != nil {
		return nil, err
	}
	return convert.MenuNodes(l.ctx, i18n.CodeByPartnerMode(l.svcCtx.Config.PartnerMode), out.GetList()), nil
}
