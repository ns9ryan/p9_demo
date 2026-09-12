package menu

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (resp *types.MenuInfo, err error) {
	out, err := l.svcCtx.Core.CreateMenu(l.ctx, convert.CreateMenuReq(req))
	if err != nil {
		return nil, err
	}
	return convert.MenuInfo(l.ctx, i18n.CodeByPartnerMode(l.svcCtx.Config.PartnerMode), out), nil
}
