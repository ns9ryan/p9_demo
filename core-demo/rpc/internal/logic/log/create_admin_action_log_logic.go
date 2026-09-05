package log

import (
	"context"

	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminActionLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAdminActionLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminActionLogLogic {
	return &CreateAdminActionLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAdminActionLogLogic) CreateAdminActionLog(in *core.CreateAdminActionLogReq) (*core.Empty, error) {
	l.svcCtx.Deps.CreateAdminActionLog(l.ctx, service.CreateAdminActionLogReq{
		UserID:         in.GetUserId(),
		RequestMethod:  in.GetRequestMethod(),
		RequestPath:    in.GetRequestPath(),
		RequestQuery:   in.GetRequestQuery(),
		RequestBody:    in.GetRequestBody(),
		ActionResult:   int16(in.GetActionResult()),
		ResponseStatus: int(in.GetResponseStatus()),
		ResponseBody:   in.GetResponseBody(),
		DurationMS:     int(in.GetDurationMs()),
		ClientIP:       in.GetClientIp(),
		UserAgent:      in.GetUserAgent(),
	})
	return &core.Empty{}, nil
}
