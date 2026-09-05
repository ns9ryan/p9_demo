package role

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

type GetRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleListLogic) GetRoleList(in *core.RoleListReq) (*core.RoleListResp, error) {
	list, total, err := l.svcCtx.Deps.ListRoles(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.RoleListReq{
		PageReq:  service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		RoleName: in.GetRoleName(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.RoleInfo, 0, len(list))
	for i := range list {
		out = append(out, logic.ToRoleInfo(&list[i]))
	}
	return &core.RoleListResp{List: out, Total: total}, nil
}
