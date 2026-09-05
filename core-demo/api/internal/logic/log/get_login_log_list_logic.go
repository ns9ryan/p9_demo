package log

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginLogListLogic {
	return &GetLoginLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginLogListLogic) GetLoginLogList(req *types.LoginLogListReq) (resp *types.LoginLogListResp, err error) {
	out, err := l.svcCtx.Core.GetLoginLogList(l.ctx, convert.LoginLogListReq(req))
	if err != nil {
		return nil, err
	}
	return convert.LoginLogList(l.ctx, out), nil
}
