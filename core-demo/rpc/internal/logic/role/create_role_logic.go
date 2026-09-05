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

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *core.CreateRoleReq) (*core.RoleInfo, error) {
	row, err := l.svcCtx.Deps.CreateRole(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), service.CreateRoleReq{
		RoleCode: in.RoleCode, RoleName: in.RoleName, Description: in.Description, Status: int16(in.Status), SortNo: int(in.SortNo),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToRoleInfo(row), nil
}
