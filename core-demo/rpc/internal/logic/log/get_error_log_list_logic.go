package log

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErrorLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetErrorLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErrorLogListLogic {
	return &GetErrorLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetErrorLogListLogic) GetErrorLogList(in *core.ErrorLogListReq) (*core.ErrorLogListResp, error) {
	list, total, err := l.svcCtx.Deps.ListErrorLogs(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.ErrorLogListReq{
		PageReq:        service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		UserID:         in.GetUserId(),
		RequestPath:    in.GetRequestPath(),
		ServiceName:    in.GetServiceName(),
		ResponseStatus: int(in.GetResponseStatus()),
		CreatedAtFrom:  in.GetCreatedAtFrom(),
		CreatedAtTo:    in.GetCreatedAtTo(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.ErrorLogInfo, 0, len(list))
	for i := range list {
		out = append(out, logic.ToErrorLogInfo(&list[i]))
	}
	return &core.ErrorLogListResp{List: out, Total: total}, nil
}
