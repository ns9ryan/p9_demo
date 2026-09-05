package log

import (
	"context"

	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateErrorLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateErrorLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateErrorLogLogic {
	return &CreateErrorLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateErrorLogLogic) CreateErrorLog(in *core.CreateErrorLogReq) (*core.Empty, error) {
	l.svcCtx.Deps.CreateErrorLog(l.ctx, service.CreateErrorLogReq{
		UserID:         in.GetUserId(),
		RequestMethod:  in.GetRequestMethod(),
		RequestPath:    in.GetRequestPath(),
		RequestQuery:   in.GetRequestQuery(),
		RequestBody:    in.GetRequestBody(),
		ServiceName:    in.GetServiceName(),
		ResponseStatus: int(in.GetResponseStatus()),
		ResponseBody:   in.GetResponseBody(),
		Subject:        in.GetSubject(),
		Detail:         in.GetDetail(),
		DurationMS:     int(in.GetDurationMs()),
		ClientIP:       in.GetClientIp(),
		UserAgent:      in.GetUserAgent(),
	})
	return &core.Empty{}, nil
}
