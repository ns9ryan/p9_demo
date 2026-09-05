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

type GetAdminActionLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminActionLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminActionLogListLogic {
	return &GetAdminActionLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminActionLogListLogic) GetAdminActionLogList(in *core.AdminActionLogListReq) (*core.AdminActionLogListResp, error) {
	list, total, err := l.svcCtx.Deps.ListAdminActionLogs(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.AdminActionLogListReq{
		PageReq:       service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		UserID:        in.GetUserId(),
		Username:      in.GetUsername(),
		RequestMethod: in.GetRequestMethod(),
		RequestPath:   in.GetRequestPath(),
		ActionResult:  int16(in.GetActionResult()),
		CreatedAtFrom: in.GetCreatedAtFrom(),
		CreatedAtTo:   in.GetCreatedAtTo(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.AdminActionLogInfo, 0, len(list))
	for i := range list {
		out = append(out, logic.ToAdminActionLogInfo(&list[i]))
	}
	return &core.AdminActionLogListResp{List: out, Total: total}, nil
}
