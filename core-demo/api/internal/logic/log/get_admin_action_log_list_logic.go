package log

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminActionLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminActionLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminActionLogListLogic {
	return &GetAdminActionLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminActionLogListLogic) GetAdminActionLogList(req *types.AdminActionLogListReq) (resp *types.AdminActionLogListResp, err error) {
	out, err := l.svcCtx.Core.GetAdminActionLogList(l.ctx, convert.AdminActionLogListReq(req))
	if err != nil {
		return nil, err
	}
	return convert.AdminActionLogList(l.ctx, out), nil
}
