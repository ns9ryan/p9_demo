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

type GetLoginLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginLogListLogic {
	return &GetLoginLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLoginLogListLogic) GetLoginLogList(in *core.LoginLogListReq) (*core.LoginLogListResp, error) {
	list, total, err := l.svcCtx.Deps.ListLoginLogs(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.LoginLogListReq{
		PageReq:     service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		Username:    in.GetUsername(),
		LoginResult: int16(in.GetLoginResult()),
		UserID:      in.GetUserId(),
		LoginAtFrom: in.GetLoginAtFrom(),
		LoginAtTo:   in.GetLoginAtTo(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.LoginLogInfo, 0, len(list))
	for i := range list {
		out = append(out, logic.ToLoginLogInfo(&list[i]))
	}
	return &core.LoginLogListResp{List: out, Total: total}, nil
}
