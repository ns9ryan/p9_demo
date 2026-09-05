package log

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErrorLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetErrorLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErrorLogListLogic {
	return &GetErrorLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetErrorLogListLogic) GetErrorLogList(req *types.ErrorLogListReq) (resp *types.ErrorLogListResp, err error) {
	out, err := l.svcCtx.Core.GetErrorLogList(l.ctx, convert.ErrorLogListReq(req))
	if err != nil {
		return nil, err
	}
	return convert.ErrorLogList(out), nil
}
