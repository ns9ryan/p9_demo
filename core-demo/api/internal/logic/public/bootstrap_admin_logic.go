package public

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BootstrapAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBootstrapAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BootstrapAdminLogic {
	return &BootstrapAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BootstrapAdminLogic) BootstrapAdmin(req *types.BootstrapAdminReq) (resp *types.UserPublic, err error) {
	out, err := l.svcCtx.Core.BootstrapAdmin(l.ctx, &coreclient.BootstrapAdminReq{
		InitToken: req.InitToken, Username: req.Username, Password: req.Password, DisplayName: req.DisplayName,
	})
	if err != nil {
		return nil, err
	}
	return convert.UserPublic(out), nil
}
